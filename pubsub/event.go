package pubsub

import "github.com/cskr/pubsub"

type Event struct {
	Event   string
	Payload interface{}
}

type Bus interface {
	Subscribe(topic string) <-chan Event
	Publish(event Event, topics ...string)
}

type InMemoryBus struct {
	pubsub   *pubsub.PubSub
	capacity int
}

func NewInMemoryBus() *InMemoryBus {
	capacity := 50
	return &InMemoryBus{
		capacity: capacity,
		pubsub:   pubsub.New(capacity),
	}
}

func (b *InMemoryBus) Subscribe(topic string) <-chan Event {
	rv := make(chan Event, b.capacity)
	internal := b.pubsub.Sub(topic)

	go func() {
		for msg := range internal {
			if event, ok := msg.(Event); ok {
				rv <- event
			}
		}
	}()

	return rv
}

func (b *InMemoryBus) Publish(event Event, topics ...string) {
	b.pubsub.Pub(event, topics...)
}
