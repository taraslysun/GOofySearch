package main

import (
	"encoding/json"
	"fmt"
	"log"
	pq "manager/priority_queue"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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

func Selector() {

}

func Router(bpq []*pq.PriorityQueue) {

}

func Prioritize(links []string) {

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
