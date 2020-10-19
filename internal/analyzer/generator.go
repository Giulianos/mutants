package analyzer

import (
	"github.com/Giulianos/mutants/internal/util"
	"strings"
	"sync"
)

// MergeGens receives generators channels
// and merges them into a single channel
// See: https://blog.go-lang.org/pipelines
type MergeGen struct{
	Generators []StrandGenerator
}

func (gen MergeGen) Generate (done <-chan struct{}) <-chan Strand {
	var wg sync.WaitGroup
	out := make(chan Strand)

	output := func(c <-chan Strand) {
		defer wg.Done()

		for strand := range c {
			select {
			case out <- strand:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(gen.Generators))
	for _, gen := range gen.Generators{
		go output(gen.Generate(done))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// HorizontalGen is a generator of horizontal strands
type HorizontalGen struct {
	dna DNA
}

func (gen HorizontalGen) Generate(done <-chan struct{}) <-chan Strand {
	out := make(chan Strand)
	go func() {
		defer close(out)
		for _, row := range gen.dna {
			select {
			case out <- Strand(row):
			case <-done:
				return
			}
		}
	}()

	return out
}

// VerticalGen is a generator of vertical strands
type VerticalGen struct {
	dna DNA
}
func (gen VerticalGen) Generate(done <-chan struct{}) <-chan Strand {
	out := make(chan Strand)
	go func() {
		defer close(out)
		if len(gen.dna) == 0{
			return
		}
		for j := range gen.dna[0] {
			var b strings.Builder
			b.Grow(len(gen.dna))

			for i := range gen.dna {
				b.WriteByte(gen.dna[i][j])
			}

			select {
			case out <- Strand(b.String()):
			case <-done:
				return
			}
		}
	}()

	return out
}

// DiagonalGen is a generator of
// diagonal (top-left -> bottom-right) strands
type DiagonalGen struct {
	dna DNA
}
func (gen DiagonalGen) Generate(done <-chan struct{}) <-chan Strand {
	out := make(chan Strand)
	go func() {
		defer close(out)

		N := len(gen.dna)
		for j := -N + 1; j < N; j++ {
			var b strings.Builder
			b.Grow(2 * N)
			for i := util.MaxInt(-j, 0); i < util.MinInt(N, N-j); i++ {
				b.WriteByte(gen.dna[i][j+i])
			}
			select {
			case out <- Strand(b.String()):
			case <-done:
				return
			}
		}
	}()

	return out
}

// AntiDiagonalGen is a generator of
// anti-diagonal (top-right -> bottom-left) strands
type AntiDiagonalGen struct {
	dna DNA
}
func (gen AntiDiagonalGen) Generate(done <-chan struct{}) <-chan Strand {
	out := make(chan Strand)
	go func() {
		defer close(out)

		N := len(gen.dna)
		for j := -N + 1; j < N; j++ {
			var b strings.Builder
			b.Grow(2 * N)
			for i := util.MaxInt(j, 0); i < util.MinInt(N, N+j); i++ {
				b.WriteByte(gen.dna[i][j+N-i-1])
			}
			select {
			case out <- Strand(b.String()):
			case <-done:
				return
			}
		}
	}()

	return out
}
