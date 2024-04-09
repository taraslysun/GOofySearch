package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	pq "manager/priority_queue"
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
}

func NewTaskManager(N int, M int) *TaskManager {

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

	fmt.Println("POST Links", links)
	tm.Prioritize(links)
	tm.Router()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// -----------------------------------------------------------

func (tm *TaskManager) checkRedis() []pq.Item {
	links, err := tm.redisClient.ZPopMax(tm.ctx, tm.redisQueue).Result()
	var redisLinks []pq.Item

	if errors.Is(err, redis.Nil) {
		return nil
	}

	if len(links) == 0 {
		return []pq.Item{}
	}

	for _, link := range links {
		redisLinks = append(redisLinks, pq.Item{
			Priority: int(link.Score),
			Value:    link.Member.(string),
		})
	}
	return redisLinks

}

func (tm *TaskManager) Selector(N int) []string {
	var links []string

	if tm.bpqs[N].Len() == 0 {
		redisLink := tm.checkRedis()[0]
		fmt.Println(redisLink.Value)
		tm.bpqs[N].Push(&redisLink)
	}
	link := tm.bpqs[N].Pop().(*pq.Item).Value
	links = append(links, link)

	return links
}

func (tm *TaskManager) Router() {
	fpqsLen := len(tm.fpqs)
	for i := 0; i < fpqsLen; i++ {
		fpq := tm.fpqs[i]
		fpqLen := fpq.Len()
		for k := 0; k < fpqLen; k++ {
			link := fpq.Pop().(*pq.Item)
			var pushed bool
			for _, q := range tm.bpqs {
				// since split url looks like this : ['https:', '', 'github.com', 'taraslysun', 'GOofySearch', 'tree', 'concurrent_crawler']
				// we take 2nd element from split array
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
			log.Fatal(err)
		}
		hostname := domain.Hostname()
		timesVisited := tm.timesVisited[hostname]
		priority := calcPriority(timesVisited, depth, len(tm.fpqs))
		tm.fpqs[priority-1].Push(&pq.Item{
			Priority: priority,
			Value:    link},
		)
		tm.timesVisited[hostname]++
	}
}

func calcPriority(timesVisited int, depth int, M int) int {
	priority := M - timesVisited*(1/2) + depth*(1/2)
	return priority
}

// -----------------------------------------------------------

func main() {
	N := 5
	M := 10
	taskManager := NewTaskManager(N, M)

	r := mux.NewRouter()
	r.Handle("/links", taskManager).Methods("GET", "POST")

	log.Println("Task Manager server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

/*
to run the server, run the following command in the task_manager directory:
   	go run manager.go
to test the server, run the following command:
	curl -X POST -d '["http://example.com"]' http://localhost:8080/links
 	curl -X GET http://localhost:8080/links
*/
