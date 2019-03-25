package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/olebedev/emitter"
	"log"
	"time"
)

type TopicCategory string
type Topic string

const (
	TopicCategoryMatchmakeRequests TopicCategory = "matchmake_requests"
)

type PubSub interface {
	Publish(topic TopicCategory, payload Event) error
	Subscribe(topic Topic) <-chan Event
}

type DbPubSub struct {
	connString string
	e *emitter.Emitter
}

type Event struct {
	Subject string `json:"subject"`
	Event string `json:"event"`
	Payload map[string]interface{} `json:payload`
}

func MkTopic(category TopicCategory, subject string) Topic {
	return Topic(string(category) + ":" + subject)
}

func NewDbPubSub(connString string)  *DbPubSub {
	emitter := emitter.New(50)
	return &DbPubSub{
		connString: connString,
		e: emitter,
	}
}

func (d *DbPubSub) Start() {
	go d.listen()
}

func (d *DbPubSub) Subscribe(topic Topic) <-chan Event {
	raw := d.e.On(string(topic))
	rv := make(chan Event, 10)

	go func() {
		for rawEvent := range raw {
			if len(rawEvent.Args) > 0 {
				if event, ok := rawEvent.Args[0].(Event); ok {
					rv <- event
				} else {
					log.Printf("error casting event payload: %+v", rawEvent)
				}
			} else {
				log.Printf("received empty event: %+v", rawEvent)
			}
		}
	}()

	return rv
}

func (d *DbPubSub) listen() {
	eventChannel := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("pubsub: error on listener channel: %s", err)
		}
	}

	listener := pq.NewListener(d.connString, time.Duration(10) * time.Millisecond, time.Duration(1) * time.Second, eventChannel)

	err := listener.Listen(string(TopicCategoryMatchmakeRequests))
	if err != nil {
		panic(fmt.Sprintf("pubsub: could not listen: %s", err))
	}
	log.Println("pubsub: listening for notifications")

	select {
	case msg := <- listener.Notify:
		d.dispatch(msg.Channel, msg.Extra)
	}
}

func (d *DbPubSub) dispatch(topic string, payload string) {
	var event Event
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		log.Printf("pubsub: unable to unmarshal db event: %s", payload)
	}

	scoped := topic + ":" + event.Subject
	d.e.Emit(scoped, event)
}

