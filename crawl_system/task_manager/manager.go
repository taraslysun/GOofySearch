package main

import (
	"context"
	pq "dcs/task_manager/priority_queue"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type TaskManager struct {
	sync.Mutex

	fpqs         []*pq.PriorityQueue
	bpqs         []*pq.PriorityQueue
	timesVisited map[string]int
	redisClient  redis.Client
	ctx          context.Context
	redisQueue   string
	chkMap       map[string]bool
	chkMapSz     int64
	curMapSz     int64
}

func NewTaskManager(N int, M int, chkMapsz int64) *TaskManager {

	var fpqs []*pq.PriorityQueue
	var bpqs []*pq.PriorityQueue

	for i := 0; i < M; i++ {
		fpqs = append(fpqs, pq.NewPriorityQueue())
	}
	for i := 0; i < N; i++ {
		bpqs = append(bpqs, pq.NewPriorityQueue())
	}

	return &TaskManager{
		fpqs:         fpqs,
		bpqs:         fpqs,
		timesVisited: make(map[string]int),
		redisClient:  *newRedisClient(),
		ctx:          newCtx(),
		redisQueue:   "redisQueue",
		chkMap:       make(map[string]bool),
		chkMapSz:     chkMapsz,
		curMapSz:     0,
	}
}

func (tm *TaskManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tm.handleGetLinks(w, r)
	case "POST":
		tm.handlePostLinks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGetLinks sends links to crawler
func (tm *TaskManager) handleGetLinks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	query := r.URL.Query()
	crawlerIdStr := query.Get("CID")

	crawlerId, err := strconv.Atoi(crawlerIdStr)
	if err != nil {
		http.Error(w, "Cannot convert id to int", http.StatusBadRequest)
	}

	links := tm.Selector(crawlerId)

	for len(links) < 15 {
		newLinks, err := tm.redisClient.ZPopMax(tm.ctx, tm.redisQueue).Result()

		if errors.Is(err, redis.Nil) {

		}
		if len(newLinks) == 0 {
			break
		}

		links = append(links, newLinks[0].Member.(string))
		newLinkDomain := strings.Split(newLinks[0].Member.(string), "/")[2]
		res, err := tm.redisClient.ZPopMax(tm.ctx, tm.redisQueue, 15).Result()
		if errors.Is(err, redis.Nil) {
			log.Fatal(err)
		}
		for _, linkObj := range res {
			link := linkObj.Member.(string)
			linkDomain := strings.Split(link, "/")[2]
			if linkDomain == newLinkDomain {
				if err != nil {
					log.Fatal(err)
				}
				links = append(links, link)
				if len(links) >= 15 {
					break
				}
			}
		}
	}

	response, err := json.Marshal(links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostLinks receives links from crawler
func (tm *TaskManager) handlePostLinks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	var links []string
	err := json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, link := range links {
		tm.addToMap(link)
	}
	tm.Prioritize(links)
	tm.Router()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// -----------------------------------------------------------

func (tm *TaskManager) addToMap(link string) {
	if tm.curMapSz > tm.chkMapSz {
		for key := range tm.chkMap {
			delete(tm.chkMap, key)
		}
		tm.chkMap[link] = true
		tm.curMapSz = 1

	} else {
		if tm.chkMap[link] {
			return
		} else {
			tm.chkMap[link] = true
			tm.curMapSz++
		}
	}
}

func (tm *TaskManager) checkRedis(N int) {

	newLink, err := tm.redisClient.ZPopMax(tm.ctx, tm.redisQueue).Result()
	fmt.Println("New link: ", newLink)

	if len(newLink) == 0 {
		return
	}

	tm.bpqs[N].Push(&pq.Item{
		Priority: int(newLink[0].Score),
		Value:    newLink[0].Member.(string),
	})

	res, err := tm.redisClient.ZPopMax(tm.ctx, tm.redisQueue, 15).Result()

	if errors.Is(err, redis.Nil) {
		log.Fatal(err)
	}

	for _, link := range res {
		linkDomain := strings.Split(link.Member.(string), "/")[2]
		if linkDomain == strings.Split(newLink[0].Member.(string), "/")[2] {
			if err != nil {
				log.Fatal(err)
			}
			if tm.bpqs[N].Len() >= 15 {
				break
			}
			tm.bpqs[N].Push(&pq.Item{Value: link.Member.(string), Priority: int(link.Score)})
		}
	}
	fmt.Println()
}

func (tm *TaskManager) Selector(N int) []string {
	var links []string
	for tm.bpqs[N].Len() != 0 {
		links = append(links, tm.bpqs[N].Pop().(*pq.Item).Value)
	}
	tm.checkRedis(N)
	return links
}

func (tm *TaskManager) Router() {
	fpqsLen := len(tm.fpqs)
	cntMap := make(map[string]int)
	for i := fpqsLen - 1; i > 0; i-- {
		fpq := tm.fpqs[i]
		fpqLen := fpq.Len()
		for k := 0; k < fpqLen; k++ {
			link := fpq.Pop().(*pq.Item)

			if len(link.Value) == 0 {
				continue
			}
			cntMap[strings.Split(link.Value, "/")[2]]++
			var pushed bool
			for _, q := range tm.bpqs {
				if (*q).Len() >= 15 {
					break
				}
				if (*q).Len() == 0 {
					q.Push(link)
					pushed = true
					break
				} else if strings.Split(link.Value, "/")[2] == strings.Split((*q)[0].Value, "/")[2] {
					q.Push(link)
					pushed = true
					break
				}
			}
			if !pushed {
				if cntMap[strings.Split(link.Value, "/")[2]] >= 15 {
					continue
				}
				err := tm.redisClient.ZAdd(tm.ctx, tm.redisQueue, &redis.Z{
					Score:  float64(link.Priority),
					Member: link.Value,
				}).Err()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// Prioritize : VD -> Visited Domains, M -> amount of FQs
func (tm *TaskManager) Prioritize(links []string) {
	for _, link := range links {
		depth := len(strings.Split(link, "/"))
		domain, err := url.Parse(link)
		if err != nil {
			continue
		}
		hostname := domain.Hostname()
		timesVisited := tm.timesVisited[hostname]
		priority := calcPriority(timesVisited, depth, len(tm.fpqs))
		tm.fpqs[priority-1].Push(
			&pq.Item{
				Priority: priority,
				Value:    link},
		)
		tm.timesVisited[hostname]++
	}
}

func calcPriority(timesVisited int, depth int, M int) int {
	priority := M - (timesVisited*(1/2) + depth*(1/2))
	return priority
}

// -----------------------------------------------------------

func main() {
	N := 8
	M := 100
	L := int64(2e4)
	taskManager := NewTaskManager(N, M, L)

	for i := 1; i < N+1; i++ {
		taskManager.checkRedis(i)
	}

	r := mux.NewRouter()
	r.Handle("/links", taskManager).Methods("GET", "POST")

	log.Println("Task Manager server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
