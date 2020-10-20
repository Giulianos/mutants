package analyzer

import (
	"encoding/json"
	"github.com/Giulianos/mutants/internal/dna"
	"github.com/Giulianos/mutants/internal/stats"
	"net/http"
)

type Controller struct {
	eventPublisher EventPublisher
}

func NewController(publisher EventPublisher) Controller {
	return Controller{
		publisher,
	}
}

func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// request structure
	data := struct {
		DNA dna.DNA `json:"dna"`
	}{}

	// decode request into data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if dna is valid
	if !dna.Validate(data.DNA) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// analyze dna and answer accordingly
	if isMutant(data.DNA) {
		// publish to stats
		c.eventPublisher.PublishVerification(stats.DNAVerification{
			DNA: data.DNA,
			Result: true,
		})
		w.WriteHeader(http.StatusOK)
		return
	}

	// publish to stats
	c.eventPublisher.PublishVerification(stats.DNAVerification{
		DNA: data.DNA,
		Result: false,
	})
	w.WriteHeader(http.StatusForbidden)
}
