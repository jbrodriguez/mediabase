package helper

import (
	"apertoire.net/mediabase/message"
	"container/heap"
	"fmt"
)

type Workpool struct {
	pool Pool
	done chan *Worker

	Work chan Request
}

func NewWorkpool(workers int, depth int) *Workpool {
	pool := make(Pool, workers)
	done := make(chan *Worker)

	for i := 0; i < workers; i++ {
		pool[i] = &Worker{make(chan Request, depth), 0, 0}
		go pool[i].work(done)
	}

	// for i := 0; i < depth; i++ {
	// 	fmt.Printf("p[i=%d]=%v\n", i, pool[i])
	// }

	work := make(chan Request)

	return &Workpool{pool, done, work}
}

func (wp *Workpool) Balance() {
	for {
		select {
		case req := <-wp.Work:
			wp.dispatch(req)
		case w := <-wp.done:
			wp.completed(w)
		}
	}
}

func (wp *Workpool) dispatch(req Request) {
	fmt.Printf("dispatch request->%v\n", req)
	// Grab the least loaded worker...
	w := heap.Pop(&wp.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&wp.pool, w)
	fmt.Printf("end dispatch request->%v\n", req)
}

func (wp *Workpool) completed(w *Worker) {
	fmt.Printf("completed worker->%s\n", w)
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	fmt.Printf("pool length is %d\n", len(wp.pool))
	heap.Remove(&wp.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&wp.pool, w)
	fmt.Printf("done completed worker->%s\n", w)
}

type Request struct {
	Fn  func(arg *message.Media) *message.Media
	Arg *message.Media
	Ch  chan *message.Media
}

type Worker struct {
	requests chan Request
	pending  int
	index    int
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests
		req.Ch <- req.Fn(req.Arg)
		done <- w
	}
}

type Pool []*Worker

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Worker)
	item.index = n
	a[n] = item
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	fmt.Printf("Pop item %d\n", len(a)-1)
	n := len(a)
	item := a[n-1]
	item.index = -1
	*p = a[0 : n-1]
	return item
}

func (p Pool) Swap(i, j int) {
	// fmt.Printf("Swap(%d, %d) and pool length is %d\n", i, j, len(p))
	// fmt.Printf("p[i=%d]=%v\n", i, p[i])
	// fmt.Printf("p[j=%d]=%v\n", j, p[j])
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
	// fmt.Printf("Swap(%d, %d) and pool length is %d\n", i, j, len(p))
}
