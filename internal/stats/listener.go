package stats

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

const (
	NSQ_TOPIC   = "events"
	NSQ_CHANNEL = "stats"
)

func StartListener(lookupAddr string, eh *EventHandler, done <-chan struct{}) {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(NSQ_TOPIC, NSQ_CHANNEL, config)
	if err != nil {
		log.Fatal(err)
	}

	consumer.ChangeMaxInFlight(200)
	consumer.AddConcurrentHandlers(eh, 20)

	err = consumer.ConnectToNSQLookupd(lookupAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-consumer.StopChan:
			return
		case <-done:
			consumer.Stop()
		}
	}

	return
}

type EventHandler struct {
	service Service
}

func NewEventHandler(service Service) EventHandler {
	return EventHandler{
		service,
	}
}

func (h *EventHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	// decode message
	var verification DNAVerification
	err := json.Unmarshal(m.Body, &verification)
	if err != nil {
		log.Println("invalid event message: ", err)
	}

	// store
	err = h.service.NotifyVerification(verification)
	if err != nil {
		log.Println("couldn't persist verification: ", err)
	}

	return nil
}
