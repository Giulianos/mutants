package dna

import "testing"

func TestValidateDNA(t *testing.T) {
	type args struct {
		dna DNA
	}
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
			if got := Validate(tt.dna); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}