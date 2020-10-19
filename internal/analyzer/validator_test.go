package analyzer

import "testing"

func Test_validateDNA(t *testing.T) {
	tests := []struct {
		name string
		dna DNA
		want bool
	}{
		{"valid dna", DNA{"ACTG", "CCTG", "AGGC", "TGGT"}, true},
		{"valid empty dna", DNA{}, true},
		{"invalid dna with lowercase", DNA{"actg", "cctg", "aggc", "tggt"}, false},
		{"invalid dna dimensions", DNA{"actg", "cctg", "aggc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateDNA(tt.dna); got != tt.want {
				t.Errorf("validateDNA() = %v, want %v", got, tt.want)
			}
		})
	}
}