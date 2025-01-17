package analyzer

// Strand is a portion of a DNA
// in any direction
type Strand string

// StrandGenerator is a contract for
// generators of strands
type StrandGenerator interface {
	Generate(done <-chan struct{}) <-chan Strand
}
