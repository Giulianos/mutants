package analyzer

import (
	"github.com/Giulianos/mutants/internal/dna"
	"testing"
)

func TestCountRepetitions(t *testing.T) {
	tests := []struct {
		strand   Strand
		expected int
	}{
		{"ACTCATCGAGATCATG", 0},
		{"CTTTTCCCGACAATTA", 1},
		{"AGGGGCTAAAATATAG", 2},
		{"AGAAAATTTTATTGTC", 2},
		{"CTAAGGGGGCGACGAC", 2},
	}

	for _, test := range tests {
		if count := countRepetitions(test.strand, 4, 3); count != test.expected {
			t.Errorf("expected %d reps, got %d", test.expected, count)
		}
	}
}

func TestIsMutant(t *testing.T) {
	tests := []struct {
		dna      dna.DNA
		expected bool
	}{
		{dna.DNA{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}, true},
		{dna.DNA{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}, false},
		{dna.DNA{}, false},
	}

	for i, test := range tests {
		if isMutant(test.dna) != test.expected {
			t.Errorf("failed test %d", i)
		}
	}
}
