package analyzer

import (
	"bytes"
	"encoding/json"
	"github.com/Giulianos/mutants/internal/dna"
	"github.com/Giulianos/mutants/internal/stats"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DummyEPMock struct {
	lastVerification stats.DNAVerification
}
func (ep *DummyEPMock) PublishVerification(verification stats.DNAVerification) {
	ep.lastVerification = verification
}

func TestController_OK(t *testing.T) {
	// mock payload
	data, err := json.Marshal(struct {
		DNA dna.DNA `json:"dna"`
	}{
		dna.DNA{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.NewBuffer(data)

	// create mock request
	req, err := http.NewRequest("GET", "/mutant", b)
	if err != nil {
		t.Fatal(err)
	}

	// mock event publisher
	ep := &DummyEPMock{}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(ep))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestController_InvalidDNA(t *testing.T) {
	// mock payload
	data, err := json.Marshal(struct {
		DNA dna.DNA `json:"dna"`
	}{
		dna.DNA{"AAAA", "AA"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.NewBuffer(data)

	// create mock request
	req, err := http.NewRequest("GET", "/mutant", b)
	if err != nil {
		t.Fatal(err)
	}

	// mock event publisher
	ep := &DummyEPMock{}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(ep))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestController_Forbidden(t *testing.T) {
	// mock payload
	data, err := json.Marshal(struct {
		DNA dna.DNA `json:"dna"`
	}{
		dna.DNA{"AA", "AA"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.NewBuffer(data)

	// create mock request
	req, err := http.NewRequest("GET", "/mutant", b)
	if err != nil {
		t.Fatal(err)
	}

	// mock event publisher
	ep := &DummyEPMock{}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(ep))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}
