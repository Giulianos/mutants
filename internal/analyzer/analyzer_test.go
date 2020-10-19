package analyzer

import "testing"

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
		if count := CountRepetitions(test.strand, 4, 3); count != test.expected {
			t.Errorf("expected %d reps, got %d", test.expected, count)
		}
	}
}

func TestIsMutant(t *testing.T) {
	tests := []struct{
		dna DNA
		expected bool
	}{
		{DNA{"ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"}, true },
		{DNA{"ATGCGA","CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}, false},
	}

	for i, test := range tests {
		if IsMutant(test.dna) != test.expected {
			t.Errorf("failed test %d", i)
		}
	}
}