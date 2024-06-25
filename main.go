package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type EnqueueRequest struct {
	Event string `json:"event"`
}
type EnqueueResponse struct {
	Message string `json:"message"`
	Event   string `json:"event"`
}

func enqueueHandler(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}

		var req EnqueueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("failed to decode request: %v", err)
			return
		}

		if err := rb.Enqueue(req.Event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("failed to enqueue item: %v", err)
			return
		}

		response := EnqueueResponse{
			Message: "Successfully enqueued event",
			Event:   req.Event,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type DequeueResponse struct {
	Message string `json:"message"`
	Event   string `json:"event"`
}

func dequeueHandler(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}

		item, err := rb.Dequeue()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("failed to dequeue item: %v", err)
			return
		}

		response := DequeueResponse{
			Message: "Successfully dequeued event",
			Event:   item.(string),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type StatusResponse struct {
	Status   string `json:"status"`
	Size     int    `json:"size"`
	Capacity int    `json:"capacity"`
	IsEmpty  bool   `json:"isEmpty"`
	IsFull   bool   `json:"isFull"`
}

func handleStatus(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}
		response := StatusResponse{
			Status:   "OK",
			Size:     rb.Size(),
			Capacity: rb.Capacity(),
			IsEmpty:  rb.IsEmpty(),
			IsFull:   rb.IsFull(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type SizeResponse struct {
	Size int `json:"size"`
}

func handleSize(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}
		response := SizeResponse{
			Size: rb.Size(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type CapacityResponse struct {
	Capacity int `json:"capacity"`
}

func handleCapacity(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}
		response := CapacityResponse{
			Capacity: rb.Capacity(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type IsEmptyResponse struct {
	IsEmpty bool `json:"isEmpty"`
}

func handleIsEmpty(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}
		response := IsEmptyResponse{
			IsEmpty: rb.IsEmpty(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

type IsFullResponse struct {
	IsFull bool `json:"isFull"`
}

func handleIsFull(rb *RingBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("unexpected method %s on endpoint %s", r.Method, r.URL.Path)
			return
		}
		response := IsFullResponse{
			IsFull: rb.IsFull(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}

func setupServer() *http.Server {
	size := os.Getenv("RING_BUFFER_SIZE")
	if size == "" {
		size = "1024"
		log.Printf("RING_BUFFER_SIZE environment variable not set, defaulting to %s", size)
	}

	log.Println("Making a new ring buffer with size:", size)

	size_int, err := strconv.Atoi(size)
	if err != nil {
		log.Fatalf("failed to convert RING_BUFFER_SIZE to int: %v", err)
	}

	rb := NewRingBuffer(size_int)
	http.HandleFunc("/enqueue", enqueueHandler(rb))
	http.HandleFunc("/dequeue", dequeueHandler(rb))
	http.HandleFunc("/status", handleStatus(rb))
	http.HandleFunc("/size", handleSize(rb))
	http.HandleFunc("/capacity", handleCapacity(rb))
	http.HandleFunc("/isEmpty", handleIsEmpty(rb))
	http.HandleFunc("/isFull", handleIsFull(rb))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT environment variable not set, defaulting to %s", port)
	}

	log.Printf("Listening at http://localhost:%s", port)

	return &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}
}

// main is excluded from test coverage
func main() {
	server := setupServer()
	if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
