package stats

import (
	"encoding/json"
	"log"
	"net/http"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service,
	}
}

func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// query service
	mutants, noMutants, ratio, err := c.service.CountVerifications()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		CountMutantDNA int64   `json:"count_mutant_dna"`
		CountHumanDNA  int64   `json:"count_human_dna"`
		Ratio          float64 `json:"ratio"`
	}{
		mutants,
		noMutants,
		ratio,
	})
}