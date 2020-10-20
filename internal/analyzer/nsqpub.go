package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Giulianos/mutants/internal/stats"
	"log"
	"net/http"
)

type NSQEventPublisher struct {
	nsqEndpoint string
}

func NewNSQEventPublisher(nsqEndpoint string) NSQEventPublisher {
	return NSQEventPublisher{
		nsqEndpoint,
	}
}

func (ep NSQEventPublisher) getTopicPublishEndpoint(topic string) string {
	return fmt.Sprintf("http://%s/pub?topic=%s", ep.nsqEndpoint, topic)
}

func (ep NSQEventPublisher) PublishVerification(verification stats.DNAVerification) {
	data, _ := json.Marshal(verification)
	buff := bytes.NewBuffer(data)
	r, err := http.Post(
		ep.getTopicPublishEndpoint("events"),
		"application/json",
		buff,
	)

	if err != nil {
		log.Println(err)
	}
	if r.StatusCode/100 != 2 {
		log.Println("error publishing message")
	}
}
