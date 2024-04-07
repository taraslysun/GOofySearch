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
	Queue []*Item
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{Queue: make([]*Item, 0)}
}

func (pq *PriorityQueue) Push(link Item) {
	pq.Lock()
	defer pq.Unlock()

	pq.Queue = append(pq.Queue, &link)
}

func (pq *PriorityQueue) Pop() Item {
	pq.Lock()
	defer pq.Unlock()

	if len(pq.Queue) == 0 {
		return Item{}
	}

	highestPriorityIdx := 0
	for i := 1; i < len(pq.Queue); i++ {
		if pq.Queue[i].Priority > pq.Queue[highestPriorityIdx].Priority {
			highestPriorityIdx = i
		}
	}
	pq.Queue[0], pq.Queue[highestPriorityIdx] = pq.Queue[highestPriorityIdx], pq.Queue[0]

	link := pq.Queue[0]
	pq.Queue = pq.Queue[1:]
	return *link
}

func (pq *PriorityQueue) IsEmpty() bool {
	pq.Lock()
	defer pq.Unlock()

	return len(pq.Queue) == 0
}

func (pq *PriorityQueue) Size() int {
	pq.Lock()
	defer pq.Unlock()

	return len(pq.Queue)
}

func (pq *PriorityQueue) Clear() {
	pq.Lock()
	defer pq.Unlock()

	pq.Queue = make([]*Item, 0)
}

func (pq *PriorityQueue) Find(value string) *Item {
	pq.Lock()
	defer pq.Unlock()

	for _, item := range pq.Queue {
		if item.Value == value {
			return item
		}
	}
	return nil
}
