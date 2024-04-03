package main

import (
	"encoding/json"
	"fmt"
	"log"
	priority_queue "manager/priority_queue"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type TaskManager struct {
	sync.Mutex
	PriorityQueue *priority_queue.PriorityQueue
	visitedLinks  map[string]bool
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		PriorityQueue: priority_queue.NewPriorityQueue(),
		visitedLinks:  make(map[string]bool),
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

func (tm *TaskManager) handleGetLinks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	var links []string
	for i := 0; i < 10; i++ {
		if tm.PriorityQueue.Size() > 0 {
			link := tm.PriorityQueue.Pop()
			tm.visitedLinks[link] = true
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
	w.Write(response)
}

func (tm *TaskManager) handlePostLinks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	var links []string
	err := json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("links post:", links)
	for _, link := range links {
		if !tm.visitedLinks[link] {
			tm.PriorityQueue.Push(link)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusCreated)
}

func main() {
	taskManager := NewTaskManager()

	r := mux.NewRouter()
	r.Handle("/links", taskManager).Methods("GET", "POST")

	log.Println("Task Manager server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
