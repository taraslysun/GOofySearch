package priority_queue

import "container/heap"

type Item struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(pq)
	heap.Fix(pq, item.index)
}

func (pq PriorityQueue) Pop() interface{} {
	old := pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	heap.Fix(pq, 0)
	return item
}
