package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	size := 3
	rb := NewRingBuffer(size)

	assert.Equal(t, size, rb.Capacity()) // the capacity of the ring buffer should be equal to the size passed to the constructor

	assert.Equal(t, 0, rb.Size()) // a new ring buffer should have a size of 0
	assert.True(t, rb.IsEmpty())  // a new ring buffer should be empty
	assert.False(t, rb.IsFull())  // a new ring buffer should not be full

	err := rb.Enqueue(1)
	assert.NoError(t, err)        // enqueuing an item should not return an error
	assert.Equal(t, 1, rb.Size()) // the size of the ring buffer should be 1 after enqueuing 1 item
	assert.False(t, rb.IsEmpty()) // the ring buffer should not be empty after enqueuing 1 item
	assert.False(t, rb.IsFull())  // the ring buffer should not be full after enqueuing 1 item

	item, err := rb.Dequeue()
	assert.NoError(t, err)        // dequeuing an item from a non-empty ring buffer should not return an error
	assert.Equal(t, 1, item)      // the dequeued item should be the same as the enqueued item
	assert.Equal(t, 0, rb.Size()) // the size of the ring buffer should be 0 after dequeuing 1 item
	assert.True(t, rb.IsEmpty())  // the ring buffer should be empty after dequeuing 1 item
	assert.False(t, rb.IsFull())  // the ring buffer should not be full after dequeuing 1 item

	// fill the ring buffer
	for i := 0; i < size; i++ {
		err := rb.Enqueue(i)
		assert.NoError(t, err) // enqueuing an item to a non-full ring buffer should not return an error
	}
	assert.Equal(t, size, rb.Size()) // the size of the ring buffer should be equal to its capacity after filling it
	assert.False(t, rb.IsEmpty())    // the ring buffer should not be empty after filling it
	assert.True(t, rb.IsFull())      // the ring buffer should be full after filling it

	// try to enqueue an item to a full ring buffer
	err = rb.Enqueue(1)
	assert.Error(t, err) // enqueuing an item to a full ring buffer should return an error

	// dequeue all items from the ring buffer
	for i := 0; i < size; i++ {
		_, err := rb.Dequeue()
		assert.NoError(t, err) // dequeuing an item from a non-empty ring buffer should not return an error
	}
	assert.Equal(t, 0, rb.Size()) // the size of the ring buffer should be 0 after dequeuing all items
	assert.True(t, rb.IsEmpty())  // the ring buffer should be empty after dequeuing all items
	assert.False(t, rb.IsFull())  // the ring buffer should not be full after dequeuing all items

	// try to dequeue an item from an empty ring buffer
	_, err = rb.Dequeue()
	assert.Error(t, err) // dequeuing an item from an empty ring buffer should return an error
}
