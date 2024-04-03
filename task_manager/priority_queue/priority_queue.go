package priority_queue

import (
	"sync"
)

type Item struct {
	Value    string
	Priority int
	index    int
}

type PriorityQueue struct {
	sync.Mutex
	queue []*Item
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{queue: make([]*Item, 0)}
}

func (pq *PriorityQueue) Push(link Item) {
	pq.Lock()
	defer pq.Unlock()

	pq.queue = append(pq.queue, &link)
}

func (pq *PriorityQueue) Pop() Item {
	pq.Lock()
	defer pq.Unlock()

	if len(pq.queue) == 0 {
		return Item{}
	}

	link := pq.queue[0]
	pq.queue = pq.queue[1:]
	return *link
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

	pq.queue = make([]*Item, 0)
}
