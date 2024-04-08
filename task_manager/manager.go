package main

import (
	"context"
	"encoding/json"
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

	links := tm.Selector("tmpQ", crawlerId)
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

	tm.Prioritize(links)
	tm.Router("tmpQ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// -----------------------------------------------------------

func (tm *TaskManager) Selector(queueName string, N int) []string {
	link, _ := tm.redisClient.LPop(tm.ctx, queueName).Result()
	links := []string{}

	if link == "" {
		for i := 0; i < 15; i++ {
			if tm.bpqs[N].Len() == 0 {
				break
			}
			link = tm.bpqs[N].Pop().(*pq.Item).Value
			links = append(links, link)
		}
	}
	return links
}

func (tm *TaskManager) Router(queueName string) {
	flen := len(tm.fpqs)
	for i := 0; i < flen; i++ {
		fmt.Println("fpqs len: ", len(tm.fpqs))
		fpq := tm.fpqs[i]
		sfpq := fpq.Len()
		for k := 0; k < sfpq; k++ {
			link := fpq.Pop().(*pq.Item)
			var pushed bool
			fmt.Println("AAAA", link.Value)
			for _, q := range tm.bpqs {
				// since split url looks like this : ['https:', '', 'github.com', 'taraslysun', 'GOofySearch', 'tree', 'concurrent_crawler']
				// we take 2nd element from split array
				if (*q).Len() == 0 {
					q.Push(link)
					fmt.Println(link.Value)
					pushed = true
					break
				}
				if strings.Split(link.Value, "/")[2] == strings.Split((*q)[0].Value, "/")[2] {
					q.Push(link)
					fmt.Println("Pushed")
					pushed = true
					break
				}
			}
			if !pushed {
				err := tm.redisClient.LPush(tm.ctx, queueName, link.Value).Err()
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
		tm.fpqs[priority-1].Push(&pq.Item{Priority: priority, Value: link})
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
