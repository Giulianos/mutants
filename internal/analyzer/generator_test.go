package analyzer

import (
	"github.com/Giulianos/mutants/internal/dna"
	"testing"
)

func TestHorizontalGen(t *testing.T) {
	tests := []struct {
		dna     dna.DNA
		strands []Strand
	}{
		{dna.DNA{"AA", "CC"}, []Strand{"AA", "CC"}},
		{dna.DNA{"ACT", "TGC", "TTC"}, []Strand{"ACT", "TGC", "TTC"}},
		{dna.DNA{}, []Strand{}},
	}

	for _, test := range tests {
		gen := horizontalGen{test.dna}
		done := make(chan struct{})
		c := gen.Generate(done)
		for _, expected := range test.strands {
			strand := <-c
			if expected != strand {
				t.Errorf("expected %s, got %s", expected, strand)
			}
		}
		close(done)
	}

}

func TestVerticalGen(t *testing.T) {
	tests := []struct {
		dna     dna.DNA
		strands []Strand
	}{
		{dna.DNA{"AA", "CC"}, []Strand{"AC", "AC"}},
		{dna.DNA{"ACT", "TGC", "TTC"}, []Strand{"ATT", "CGT", "TCC"}},
		{dna.DNA{}, []Strand{}},
	}

	for _, test := range tests {
		gen := verticalGen{test.dna}
		done := make(chan struct{})
		c := gen.Generate(done)
		for _, expected := range test.strands {
			strand := <-c
			if expected != strand {
				t.Errorf("expected %s, got %s", expected, strand)
			}
		}
		close(done)
	}

}

func TestDiagonalGen(t *testing.T) {
	tests := []struct {
		dna     dna.DNA
		strands []Strand
	}{
		{dna.DNA{"AT", "CC"}, []Strand{"C", "AC", "T"}},
		{dna.DNA{"ACT", "TGC", "TTC"}, []Strand{"T", "TT", "AGC", "CC", "T"}},
		{dna.DNA{}, []Strand{}},
	}

	for _, test := range tests {
		gen := diagonalGen{test.dna}
		done := make(chan struct{})
		c := gen.Generate(done)
		for _, expected := range test.strands {
			strand := <-c
			if expected != strand {
				t.Errorf("expected %s, got %s", expected, strand)
			}
		}
		close(done)
	}

}

func TestAntiDiagonalGen(t *testing.T) {
	tests := []struct {
		dna     dna.DNA
		strands []Strand
	}{
		{dna.DNA{"AT", "CC"}, []Strand{"A", "TC", "C"}},
		{dna.DNA{"ACT", "TGC", "TTC"}, []Strand{"A", "CT", "TGT", "CT", "C"}},
		{dna.DNA{}, []Strand{}},
	}

	for _, test := range tests {
		gen := antiDiagonalGen{test.dna}
		done := make(chan struct{})
		c := gen.Generate(done)
		for _, expected := range test.strands {
			strand := <-c
			if expected != strand {
				t.Errorf("expected %s, got %s", expected, strand)
			}
		}
		close(done)
	}

}

func createStrandSet(generators []StrandGenerator) (int, map[Strand]bool) {
	set := map[Strand]bool{}
	var count int
	for _, gen := range generators {
		done := make(chan struct{})
		c := gen.Generate(done)
		for strand := range c {
			set[strand] = true
			count++
		}
		close(done)
	}

	return count, set
}

func TestMergeGen(t *testing.T) {
	tests := []dna.DNA{
		{"AT", "CC"},
		{"ACT", "TGC", "TTC"},
		{},
	}

	for _, dnaSeq := range tests {
		// create generators
		gens := []StrandGenerator{
			verticalGen{dnaSeq},
			horizontalGen{dnaSeq},
			diagonalGen{dnaSeq},
			antiDiagonalGen{dnaSeq},
		}

		// build expected set
		expectedCount, expectedValues := createStrandSet(gens)

		// merge generators
		done := make(chan struct{})
		c := mergeGen{gens}.Generate(done)

		var count int
		for strand := range c {
			if _, exists := expectedValues[strand]; !exists {
				t.Error("unexpected strand:", strand)
			} else {
				count++
			}
		}

		if count != expectedCount {
			t.Errorf("merged generator must generate all strands, expected: %d, got %d", count, expectedCount)
		}
	}
}

func TestExplicitCancellation(t *testing.T) {
	dnaSeq := dna.DNA{
		"ACTTAGTG",
		"CATAGTTA",
		"ATGCAATG",
		"TTGACCCT",
		"TTGGCATG",
		"ATTCGAAT",
		"ATTGCCAA",
		"GTTACGGT",
	}

	gens := []StrandGenerator{
		verticalGen{dnaSeq},
		horizontalGen{dnaSeq},
		diagonalGen{dnaSeq},
		antiDiagonalGen{dnaSeq},
	}

	done := make(chan struct{})
	c := mergeGen{gens}.Generate(done)

	// cancell generator
	close(done)

	var count int
	for _ = range c {
		count++
	}

	// generated strands after cancellation should be less
	// or equal than the quantity of generators
	if count > len(gens) {
		t.Errorf("generators should stop after cancellation")
	}
}
