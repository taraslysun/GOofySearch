package priority_queue

import "container/heap"

type Item struct {
	Value    string
	Priority int
	Index    int
}

type PriorityQueue []*Item

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

func (pq PriorityQueue) isEmpty() bool {
	return len(pq) == 0
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.Index = len(*pq)
	*pq = append(*pq, item)
	heap.Fix(pq, item.Index)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	heap.Fix(pq, 0)
	return item
}
