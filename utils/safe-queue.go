package utils

import "sync"

type SafeQueue struct {
	items []interface{}
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewSafeQueue() *SafeQueue {
	sq := &SafeQueue{}
	sq.cond = sync.NewCond(&sq.mu)
	return sq
}

func (sq *SafeQueue) Enqueue(item interface{}) {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	sq.items = append(sq.items, item)
	sq.cond.Signal()
}

func (sq *SafeQueue) Dequeue() interface{} {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	for len(sq.items) == 0 {
		sq.cond.Wait()
	}

	item := sq.items[0]
	sq.items = sq.items[1:]
	return item
}

func (sq *SafeQueue) Len() int {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	return len(sq.items)
}

func (sq *SafeQueue) Peek() interface{} {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	if len(sq.items) == 0 {
		return nil
	}
	return sq.items[0]
}
