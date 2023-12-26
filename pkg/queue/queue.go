package queue

import (
	"fmt"

	"github.com/yash492/statusy/pkg/schema"
)

type IncidentPayload struct {
	IncidentUpdate schema.IncidentUpdate
	// State can be open/updated/closed
	State string
}

type Queue struct {
	channel chan IncidentPayload
}

var ErrQueueFull = fmt.Errorf("queue is full")
var ErrQueueEmpty = fmt.Errorf("queue is empty")

func New(capacity int) *Queue {
	return &Queue{
		channel: make(chan IncidentPayload, capacity),
	}
}

func (q *Queue) Publish(value IncidentPayload) error {
	select {
	case q.channel <- value:
		return nil
	default:
		return ErrQueueFull

	}
}

func (q *Queue) Consume() (IncidentPayload, error) {
	select {
	case incident := <-q.channel:
		return incident, nil
	default:
		return IncidentPayload{}, ErrQueueEmpty
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.channel) < 1
}
