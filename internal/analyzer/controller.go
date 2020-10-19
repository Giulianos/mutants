package analyzer

import (
	"encoding/json"
	"net/http"
)

func PostController(w http.ResponseWriter, r *http.Request) {
	// request structure
	data := struct {
		DNA DNA `json:"dna"`
	}{}

	// decode request into data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if dna is valid
	if !validateDNA(data.DNA) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// analyze dna and answer accordingly
	if isMutant(data.DNA) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusForbidden)
}
