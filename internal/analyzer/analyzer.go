package analyzer

// CountRepetitions counts, at most, limit repetitions
// of length nitrogenous bases in strand.
func CountRepetitions(strand Strand, length, limit int) int {
	currLen := 0
	reps := 0
	var prev rune
	for _, r := range strand {
		if prev != r {
			currLen = 0
		}
		currLen++
		if  currLen >= length {
			reps++ // sequences can overlap
		}
		if limit > 0 && reps >= limit {
			break
		}
		prev = r
	}

	return reps
}

// IsMutant analyzes from a dna sample,
// whether its or not a mutant
func IsMutant(dna DNA) bool {
	gens := []StrandGenerator{
		VerticalGen{dna},
		HorizontalGen{dna},
		DiagonalGen{dna},
		AntiDiagonalGen{dna},
	}

	done := make(chan struct{})
	defer close(done)

	c := MergeGen{gens}.Generate(done)

	var repsCount int
	for strand := range c {
		repsCount += CountRepetitions(strand, 4, 2)
		if repsCount >= 2 {
			return true
		}
	}

	return false
}