package analyzer

import "github.com/Giulianos/mutants/internal/dna"

// countRepetitions counts, at most, limit repetitions
// of length nitrogenous bases in strand.
func countRepetitions(strand Strand, length, limit int) int {
	currLen := 0
	reps := 0
	var prev rune
	for _, r := range strand {
		if prev != r {
			currLen = 0
		}
		currLen++
		if currLen >= length {
			reps++ // sequences can overlap
		}
		if limit > 0 && reps >= limit {
			break
		}
		prev = r
	}

	return reps
}

// isMutant analyzes from a dna sample,
// whether its or not a mutant
func isMutant(dna dna.DNA) bool {
	gens := []StrandGenerator{
		verticalGen{dna},
		horizontalGen{dna},
		diagonalGen{dna},
		antiDiagonalGen{dna},
	}

	done := make(chan struct{})
	defer close(done)

	c := mergeGen{gens}.Generate(done)

	var repsCount int
	for strand := range c {
		repsCount += countRepetitions(strand, 4, 2)
		if repsCount >= 2 {
			return true
		}
	}

	return false
}
