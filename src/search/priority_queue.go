package search

type priorityQueue struct {
	items []*Node
	priority func(*Node) int
}

func newPriorityQueue(priority func(*Node) int) *priorityQueue {
	return &priorityQueue{priority: priority}
}

func (pq *priorityQueue) Len() int { 
	return len(pq.items)
}

func (pq *priorityQueue) Less(i, j int) bool {
	return pq.priority(pq.items[i]) < pq.priority(pq.items[j])
}

func (pq *priorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *priorityQueue) Push(x any) {
	pq.items = append(pq.items, x.(*Node))
}

func (pq *priorityQueue) Pop() any {
	n := len(pq.items)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	return item
}
