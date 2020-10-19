package analyzer

import (
	"fmt"
	"regexp"
)

func validateDNA(dna DNA) bool {
	N := len(dna)
	re, _ := regexp.Compile(fmt.Sprintf("[ACTG]{%d}", N))

	for _, s := range dna {
		if !re.MatchString(s) {
			return false
		}
	}

	return true
}
