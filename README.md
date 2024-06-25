[![Coverage Status](https://img.shields.io/badge/coverage-80.9%25-brightgreen)](https://github.com/username/repo)

# Readme

A basic event queue service implemented in Go. Exposes a REST API to add events to a queue and retrieve events from the queue.

## What is an Event?

An event is a stringified JSON object. It can be any arbitrary JSON object. The event queue service does not check the structure of the event, it just stores it and delivers it as a string.

## API

### Enqueue Event

Adds an event to the queue.

#### Request

```http
POST /enqueue
Content-Type: application/json

{
  "event": "{\n  \"name\": \"EVENT_NAME\",\n  \"value\": 42\n}"
}
```

#### Response

```http
HTTP/1.1 201 Created
Content-Type: application/json

{
  "message": "Successfully enqueued event",
  "event": "123"
}
```

### Dequeue Event

Retrieves the next event from the queue.

#### Request

```http
GET /dequeue
```

#### Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "event": "{\n  \"name\": \"EVENT_NAME\",\n  \"value\": 42\n}"
}
```

### Status

Checks the status of the service.

#### Request

```http
GET /status
```

#### Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "OK"
  "size": 0,
  "capacity": 100,
  "isEmpty": true,
  "isFull": false
}
```
