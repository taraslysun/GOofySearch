package priority_queue

import (
	"sync"
)

type PriorityQueue struct {
	sync.Mutex
	queue []string
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		queue: make([]string, 0),
	}
}

func (pq *PriorityQueue) Push(link string) {
	pq.Lock()
	defer pq.Unlock()

	pq.queue = append(pq.queue, link)
}

func (pq *PriorityQueue) Pop() string {
	pq.Lock()
	defer pq.Unlock()

	if len(pq.queue) == 0 {
		return ""
	}

	link := pq.queue[0]
	pq.queue = pq.queue[1:]
	return link
}

func (pq *PriorityQueue) IsEmpty() bool {
	pq.Lock()
	defer pq.Unlock()

	return len(pq.queue) == 0
}

func (pq *PriorityQueue) Size() int {
	pq.Lock()
	defer pq.Unlock()

	return len(pq.queue)
}

func (pq *PriorityQueue) Clear() {
	pq.Lock()
	defer pq.Unlock()

	pq.queue = make([]string, 0)
}
