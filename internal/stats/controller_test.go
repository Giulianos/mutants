package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepo struct {
	humans int64
	mutants int64
	err error
}

func (mr *mockRepo) Persist(verification DNAVerification) error {
	mr.humans++
	if verification.Result {
		mr.mutants++
	}

	return mr.err
}

func (mr *mockRepo) CountByResult(value bool) (int64, error) {
	if value {
		return mr.mutants, mr.err
	} else {
		return mr.humans - mr.mutants, mr.err
	}
}

func TestController_CorrectValues(t *testing.T) {
	// create mock request
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock repo
	repo := &mockRepo{
		humans: 10,
		mutants: 1,
	}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(NewService(repo)))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	// check response
	data := struct {
		CountMutantDNA int64 `json:"count_mutant_dna"`
		CountHumanDNA int64 `json:"count_human_dna"`
		Ratio float64 `json:"ratio"`
	}{}
	json.NewDecoder(rr.Body).Decode(&data)
	if data.CountHumanDNA != 10 {
		t.Errorf("handler returned wrong count_human_dna: got %v want %v", data.CountHumanDNA,
		10)
	}
	if data.CountMutantDNA != 1 {
		t.Errorf("handler returned wrong count_mutant_dna: got %v want %v", data.CountMutantDNA,
			1)
	}
	if data.Ratio != 0.1 {
		t.Errorf("handler returned wrong ratio: got %v want %v", data.Ratio,
			0.1)
	}
}


func TestController_ServerError(t *testing.T) {
	// create mock request
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock repo
	repo := &mockRepo{
		humans: 10,
		mutants: 1,
		err: fmt.Errorf("unexpected db error"),
	}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(NewService(repo)))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusInternalServerError{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}


func TestController_EmptyRepo(t *testing.T) {
	// create mock request
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock repo
	repo := &mockRepo{
		humans: 0,
		mutants: 0,
	}

	// handle request, record response
	rr := httptest.NewRecorder()
	handler := http.Handler(NewController(NewService(repo)))
	handler.ServeHTTP(rr, req)

	// check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	// check response
	data := struct {
		CountMutantDNA int64 `json:"count_mutant_dna"`
		CountHumanDNA int64 `json:"count_human_dna"`
		Ratio float64 `json:"ratio"`
	}{}
	json.NewDecoder(rr.Body).Decode(&data)
	if data.CountHumanDNA != 0 {
		t.Errorf("handler returned wrong count_human_dna: got %v want %v", data.CountHumanDNA,
			0)
	}
	if data.CountMutantDNA != 0 {
		t.Errorf("handler returned wrong count_mutant_dna: got %v want %v", data.CountMutantDNA,
			0)
	}
	if data.Ratio != 0 {
		t.Errorf("handler returned wrong ratio: got %v want %v", data.Ratio,
			0)
	}
}
