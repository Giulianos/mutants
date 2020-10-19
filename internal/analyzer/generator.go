package analyzer

import (
	"github.com/Giulianos/mutants/internal/util"
	"strings"
	"sync"
)

// mergeGen receives generators channels
// and merges them into a single channel
// See: https://blog.go-lang.org/pipelines
type mergeGen struct {
	Generators []StrandGenerator
}

func (gen mergeGen) Generate(done <-chan struct{}) <-chan Strand {
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
	for _, gen := range gen.Generators {
		go output(gen.Generate(done))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// horizontalGen is a generator of horizontal strands
type horizontalGen struct {
	dna DNA
}

func (gen horizontalGen) Generate(done <-chan struct{}) <-chan Strand {
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

// verticalGen is a generator of vertical strands
type verticalGen struct {
	dna DNA
}

func (gen verticalGen) Generate(done <-chan struct{}) <-chan Strand {
	out := make(chan Strand)
	go func() {
		defer close(out)
		if len(gen.dna) == 0 {
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

// diagonalGen is a generator of
// diagonal (top-left -> bottom-right) strands
type diagonalGen struct {
	dna DNA
}

func (gen diagonalGen) Generate(done <-chan struct{}) <-chan Strand {
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

// antiDiagonalGen is a generator of
// anti-diagonal (top-right -> bottom-left) strands
type antiDiagonalGen struct {
	dna DNA
}

func (gen antiDiagonalGen) Generate(done <-chan struct{}) <-chan Strand {
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
