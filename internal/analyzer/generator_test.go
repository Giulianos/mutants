package analyzer

import (
	"testing"
)

func TestHorizontalGen(t *testing.T) {
	tests := []struct {
		dna     DNA
		strands []Strand
	}{
		{DNA{"AA", "CC"}, []Strand{"AA", "CC"}},
		{DNA{"ACT", "TGC", "TTC"}, []Strand{"ACT", "TGC", "TTC"}},
		{DNA{}, []Strand{}},
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
		dna     DNA
		strands []Strand
	}{
		{DNA{"AA", "CC"}, []Strand{"AC", "AC"}},
		{DNA{"ACT", "TGC", "TTC"}, []Strand{"ATT", "CGT", "TCC"}},
		{DNA{}, []Strand{}},
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
		dna     DNA
		strands []Strand
	}{
		{DNA{"AT", "CC"}, []Strand{"C", "AC", "T"}},
		{DNA{"ACT", "TGC", "TTC"}, []Strand{"T", "TT", "AGC", "CC", "T"}},
		{DNA{}, []Strand{}},
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
		dna     DNA
		strands []Strand
	}{
		{DNA{"AT", "CC"}, []Strand{"A", "TC", "C"}},
		{DNA{"ACT", "TGC", "TTC"}, []Strand{"A", "CT", "TGT", "CT", "C"}},
		{DNA{}, []Strand{}},
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
	tests := []DNA{
		{"AT", "CC"},
		{"ACT", "TGC", "TTC"},
		{},
	}

	for _, dna := range tests {
		// create generators
		gens := []StrandGenerator{
			verticalGen{dna},
			horizontalGen{dna},
			diagonalGen{dna},
			antiDiagonalGen{dna},
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
	dna := DNA{
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
		verticalGen{dna},
		horizontalGen{dna},
		diagonalGen{dna},
		antiDiagonalGen{dna},
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
