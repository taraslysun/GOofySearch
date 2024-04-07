package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	pq "manager/priority_queue"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type TaskManager struct {
	sync.Mutex
	PriorityQueue *pq.PriorityQueue
	visitedLinks  map[string]bool

	fpqs         []*pq.PriorityQueue
	bpqs         []*pq.PriorityQueue
	timesVisited map[string]int
	redisClient  redis.Client
	ctx          context.Context
}

func NewTaskManager(N int, M int) *TaskManager {
	return &TaskManager{
		PriorityQueue: pq.NewPriorityQueue(),
		visitedLinks:  make(map[string]bool),
		fpqs:          make([]*pq.PriorityQueue, M),
		bpqs:          make([]*pq.PriorityQueue, N),
		timesVisited:  make(map[string]int),
		redisClient:   *newRedisClient(),
		ctx:           newCtx(),
	}
}

func (tm *TaskManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tm.handleGetLinks(w)
	case "POST":
		tm.handlePostLinks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGetLinks sends links to crawler
func (tm *TaskManager) handleGetLinks(w http.ResponseWriter) {
	tm.Lock()
	defer tm.Unlock()

	links := tm.Selector("tmpQ", 3)
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

func (tm *TaskManager) Selector(queueName string, N int) string {
	link, err := tm.redisClient.LPop(tm.ctx, queueName).Result()
	if err != nil {
		log.Fatal(err)
	}
	if link == "" {
		return ""
	}
	link = tm.bpqs[N].Pop().Value
	return link
}

func (tm *TaskManager) Router(queueName string) {
	for len(tm.fpqs) != 0 {
		fpq := tm.fpqs[len(tm.fpqs)-1]
		tm.fpqs = tm.fpqs[:len(tm.fpqs)-1]
		for !fpq.IsEmpty() {
			link := fpq.Pop()
			var pushed bool
			for _, q := range tm.bpqs {
				// since split url looks like this : ['https:', '', 'github.com', 'taraslysun', 'GOofySearch', 'tree', 'concurrent_crawler']
				// we take 2nd element from split array
				if strings.Split(link.Value, "/")[2] == strings.Split(q.Queue[0].Value, "/")[2] {
					q.Push(link)
					pushed = true
					break
				}
			}
			if !pushed {
				err := tm.redisClient.LPush(tm.ctx, queueName, link).Err()
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
		tm.fpqs[priority].Push(pq.Item{Priority: priority, Value: link})
		tm.timesVisited[hostname]++
	}
}

func calcPriority(timesVisited int, depth int, M int) int {
	priority := M - timesVisited*(1/2) + depth*(1/2)
	return priority
}

// -----------------------------------------------------------

func main() {
	N := 8
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
