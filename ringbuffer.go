package main

import "errors"

type RingBuffer struct {
	buffer []interface{}
	size   int
	head   int
	tail   int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, capacity),
		size:   0,
		head:   0,
		tail:   0,
	}
}

func (rb *RingBuffer) Enqueue(item interface{}) error {
	if rb.size == len(rb.buffer) {
		return errors.New("error: cannot enqueue to a full buffer")
	}
	rb.buffer[rb.tail] = item
	rb.tail = (rb.tail + 1) % len(rb.buffer)
	rb.size++
	return nil
}

func (rb *RingBuffer) Dequeue() (interface{}, error) {
	if rb.size == 0 {
		return nil, errors.New("error: cannot dequeue from an empty buffer")
	}
	item := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % len(rb.buffer)
	rb.size--
	return item, nil
}

func (rb *RingBuffer) Size() int {
	return rb.size
}

func (rb *RingBuffer) Capacity() int {
	return len(rb.buffer)
}

func (rb *RingBuffer) IsEmpty() bool {
	return rb.size == 0
}

func (rb *RingBuffer) IsFull() bool {
	return rb.size == len(rb.buffer)
}
