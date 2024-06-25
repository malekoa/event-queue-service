package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func disableLogging() func() {
	originalOutput := log.Writer()
	log.SetOutput(io.Discard)
	return func() {
		log.SetOutput(originalOutput)
	}
}

func TestEnqueueHandler(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := enqueueHandler(rb)
		req := httptest.NewRequest(http.MethodGet, "/enqueue", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("failed to decode request body", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := enqueueHandler(rb)
		req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "invalid character 'i' looking for beginning of value\n", w.Body.String())
	})

	t.Run("failed to enqueue item", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := enqueueHandler(rb)
		req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewBufferString(`{"event": "event1"}`))
		w := httptest.NewRecorder()

		rb.Enqueue("event1")
		rb.Enqueue("event2")
		rb.Enqueue("event3")
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "error: cannot enqueue to a full buffer\n", w.Body.String())
	})

	t.Run("successfully enqueued event", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := enqueueHandler(rb)
		req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewBufferString(`{"event": "event1"}`))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestDequeueHandler(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := dequeueHandler(rb)
		req := httptest.NewRequest(http.MethodPost, "/dequeue", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully dequeued event", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := dequeueHandler(rb)
		req := httptest.NewRequest(http.MethodGet, "/dequeue", nil)
		w := httptest.NewRecorder()

		rb.Enqueue("event1")
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"message":"Successfully dequeued event","event":"event1"}`+"\n", w.Body.String())
	})

	t.Run("failed to dequeue item", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := dequeueHandler(rb)
		req := httptest.NewRequest(http.MethodGet, "/dequeue", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "error: cannot dequeue from an empty buffer\n", w.Body.String())
	})
}

func TestHandleStatus(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleStatus(rb)
		req := httptest.NewRequest(http.MethodPost, "/status", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully retrieved status", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleStatus(rb)
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSetupServer(t *testing.T) {
	defer disableLogging()()
	t.Run("setup server", func(t *testing.T) {
		server := setupServer()
		assert.NotNil(t, server)
	})
}

func TestHandlsSize(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleSize(rb)
		req := httptest.NewRequest(http.MethodPost, "/size", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully retrieved size", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleSize(rb)
		req := httptest.NewRequest(http.MethodGet, "/size", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandleCapacity(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleCapacity(rb)
		req := httptest.NewRequest(http.MethodPost, "/capacity", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully retrieved capacity", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleCapacity(rb)
		req := httptest.NewRequest(http.MethodGet, "/capacity", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandleIsEmpty(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleIsEmpty(rb)
		req := httptest.NewRequest(http.MethodPost, "/isEmpty", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully retrieved isEmpty", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleIsEmpty(rb)
		req := httptest.NewRequest(http.MethodGet, "/isEmpty", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandleIsFull(t *testing.T) {
	defer disableLogging()()

	t.Run("method not allowed", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleIsFull(rb)
		req := httptest.NewRequest(http.MethodPost, "/isFull", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "Method not allowed\n", w.Body.String())
	})

	t.Run("successfully retrieved isFull", func(t *testing.T) {
		rb := NewRingBuffer(3)
		handler := handleIsFull(rb)
		req := httptest.NewRequest(http.MethodGet, "/isFull", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
