package main

import (
	"context"
	"encoding/json"
	"fmt"
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
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		PriorityQueue: pq.NewPriorityQueue(),
		visitedLinks:  make(map[string]bool),
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

func (tm *TaskManager) handleGetLinks(w http.ResponseWriter) {
	tm.Lock()
	defer tm.Unlock()

	var links []pq.Item
	for i := 0; i < 10; i++ {
		if tm.PriorityQueue.Size() > 0 {
			link := tm.PriorityQueue.Pop()
			tm.visitedLinks[link.Value] = true
			links = append(links, link)
		} else {
			break
		}
	}

	fmt.Println("links get:", links)

	if links == nil {
		http.Error(w, "No links available", http.StatusNoContent)
		return
	}

	response, err := json.Marshal(links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func (tm *TaskManager) handlePostLinks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	var links []pq.Item
	err := json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("links post:", links)
	for _, link := range links {
		if !tm.visitedLinks[link.Value] {
			tm.PriorityQueue.Push(link)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
}

// -----------------------------------------------------------

func Selector(bpqs []*pq.PriorityQueue, redisClient redis.Client, queueName string, ctx context.Context, N int) string {
	link, err := redisClient.LPop(ctx, queueName).Result()
	if err != nil {
		log.Fatal(err)
	}
	if link != "" {
		return link
	}
	link = bpqs[N].Pop().Value
	return link
}

func Router(fpqs []*pq.PriorityQueue, bpqs []*pq.PriorityQueue, redisClient redis.Client, queueName string, ctx context.Context) {
	for len(fpqs) != 0 {
		fpq := fpqs[len(fpqs)-1]
		fpqs = fpqs[:len(fpqs)-1]
		for !fpq.IsEmpty() {
			link := fpq.Pop()
			var pushed bool
			for _, q := range bpqs {
				if strings.Split(link.Value, "/")[2] == strings.Split(q.Queue[0].Value, "/")[2] {
					q.Push(link)
					pushed = true
					break
				}
			}
			if !pushed {
				err := redisClient.LPush(ctx, queueName, link).Err()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// Prioritize : VD -> Visited Domains, M -> amount of FQs
func Prioritize(links []string, M int, VD map[string]int, fpqs []*pq.PriorityQueue) {
	for _, link := range links {
		depth := len(strings.Split(link, "/"))
		domain, err := url.Parse(link)
		if err != nil {
			return
		}
		hostname := domain.Hostname()
		timesVisited := VD[hostname]
		priority := calcPriority(timesVisited, depth, M)
		fpqs[priority].Push(pq.Item{Priority: priority, Value: link})
		VD[hostname]++
	}
}

func calcPriority(timesVisited int, depth int, M int) int {
	priority := M - timesVisited*(1/2) + depth*(1/2)
	return priority
}

// -----------------------------------------------------------

func main() {
	taskManager := NewTaskManager()

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
