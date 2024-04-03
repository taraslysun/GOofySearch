package concurrent_manager

import (
	"container/heap"
)

type Item struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority > (*pq)[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	if len(*pq) == 0 {
		return nil
	}
	old := (*pq)[0]
	(*pq)[0], (*pq)[len(*pq)-1] = (*pq)[len(*pq)-1], (*pq)[0]
	*pq = (*pq)[:len(*pq)-1]
	heap.Fix(pq, old.index)
	return old
}
