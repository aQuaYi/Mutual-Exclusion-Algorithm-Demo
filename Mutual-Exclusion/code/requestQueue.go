package main

import (
	"container/heap"
	"sync"
)

// RequestQueue 提供了操作 request queue 的接口
type RequestQueue interface {
	// Min 返回最小的 Timestamp 值
	Min() Less
	// Push 把元素加入 RequestQueue 中
	Push(Less)
	// Remove 在 RequestQueue 中删除 Less
	Remove(Less)
}

type requestQueue struct {
	rpq       *requestPriorityQueue
	requestOf map[Less]*request
	mutex     sync.Mutex
}

func newRequestQueue() RequestQueue {
	return &requestQueue{
		rpq:       new(requestPriorityQueue),
		requestOf: make(map[Less]*request, 1024),
	}
}

func (rq *requestQueue) Min() Less {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	if len(*rq.rpq) == 0 {
		return nil
	}
	return (*rq.rpq)[0].ls
}

func (rq *requestQueue) Push(ls Less) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	r := &request{
		ls: ls,
	}

	rq.requestOf[ls] = r
	heap.Push(rq.rpq, r)
}

func (rq *requestQueue) Remove(ls Less) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	rq.rpq.remove(rq.requestOf[ls])
	delete(rq.requestOf, ls)
}

// Less 是 rpq 元素中的主要成分
type Less interface {
	// Less 比较两个接口的值
	Less(interface{}) bool
}

// request 是 priorityQueue 中的元素
type request struct {
	ls    Less
	index int
}

// rpq implements heap.Interface and holds entries.
type requestPriorityQueue []*request

func (q requestPriorityQueue) Len() int { return len(q) }

func (q requestPriorityQueue) Less(i, j int) bool {
	return q[i].ls.Less(q[j].ls)
}

func (q requestPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push 往 pq 中放 entry
func (q *requestPriorityQueue) Push(x interface{}) {
	temp := x.(*request)
	temp.index = len(*q)
	*q = append(*q, temp)
}

// Pop 从 pq 中取出最优先的 entry
func (q *requestPriorityQueue) Pop() interface{} {
	temp := (*q)[len(*q)-1]
	temp.index = -1 // for safety
	*q = (*q)[0 : len(*q)-1]
	return temp
}

func (q *requestPriorityQueue) remove(r *request) {
	heap.Remove(q, r.index)
}
