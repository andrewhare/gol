package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// Pattern is a matrix of runes containing
// live and dead cells.
type Pattern [][]rune

const (
	// L is a live cell.
	L = 'O'
	// D is a dead cell.
	D = ' '
)

func main() {
	// A map of patterns that can be used as the initial
	// state of a new board.
	patterns := map[string]Pattern{
		"glider": {
			{D, L, D},
			{D, D, L},
			{L, L, L},
		},
		"spaceship": {
			{D, L, D, D, L},
			{L, D, D, D, D},
			{L, D, D, D, L},
			{L, L, L, L, D},
			{D, D, D, D, D},
		},
	}

	var (
		dimensions = flag.Int("dimensions", 25, "board dimensions")
		pattern    = flag.String("pattern", "glider", "initial pattern on the board")
		interval   = flag.Duration("interval", time.Second, "the interval between generations")
	)

	flag.Parse()

	p, ok := patterns[*pattern]
	if !ok {
		log.Fatalf("unknown pattern: %s\n", *pattern)
	}

	var (
		b = NewBoard(*dimensions, p)
		c = make(chan string)
	)

	go func(ch <-chan string) {
		area, err := pterm.DefaultArea.WithFullscreen().Start()
		if err != nil {
			log.Fatal(err)
		}

		for s := range ch {
			area.Update(s)
		}
	}(c)

	for {
		c <- b.String()
		time.Sleep(*interval)
		b.Tick()
	}
}

// NewBoard creates a new board with the given dimensions and pattern.
func NewBoard(dim int, p Pattern) *Board {
	var (
		// Assuming all patterns will be square matrices,
		// therefore the pattern dimensions are equal to the length
		// of one side of the pattern.
		pDim = len(p)
		// The pattern is going to be centered. This is the first
		// row on which to start writing pattern data.
		pFirstRow = (dim - pDim) / 2
		// The last row on which to write pattern data.
		pLastRow = pFirstRow + pDim
		// Patterns are squares, first column is the same as
		// the first row.
		pFirstCol = pFirstRow
		// The last column on which to write pattern data.
		pLastCol = pFirstCol + pDim
	)

	gen := make(Pattern, dim)
	for row := 0; row < dim; row++ {
		gen[row] = make([]rune, dim)
		for col := 0; col < dim; col++ {
			// Determine if the current row and column fall within the bounds of the pattern.
			if row >= pFirstCol && row < pLastCol && col >= pFirstRow && col < pLastRow {
				// Instead of writing a dead cell, write the value from the
				// pattern by offsetting the value from the current column and row.
				gen[row][col] = p[row-pFirstCol][col-pFirstRow]
				continue
			}

			gen[row][col] = D
		}
	}

	return &Board{
		dim: dim,
		gen: gen,
	}
}

// Board represents a game board that can be rendered
// as string output.
type Board struct {
	dim int
	gen Pattern
}

// Tick creates a new generation of the game board based on
// its current state.
func (b *Board) Tick() {
	nextGen := make(Pattern, b.dim)

	for row := 0; row < len(b.gen); row++ {
		nextGen[row] = make([]rune, b.dim)
		for col := 0; col < len(b.gen); col++ {
			// Set the fate of the current cell based
			// on the number of neighbors it has.
			nextGen[row][col] = b.fate(row, col)
		}
	}

	b.gen = nextGen
}

// String returns the current board state as a string.
func (b *Board) String() string {
	var sb strings.Builder

	for row := 0; row < len(b.gen); row++ {
		sb.WriteString(string(b.gen[row]))
		sb.WriteString("\n")
	}

	return sb.String()
}

// Fate determines the cell's fate based on the
// current generation and returns the cell's new
// state for the next generation.
func (b *Board) fate(cellRow, cellCol int) rune {
	var neighbors int
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			var (
				// The column and row of the neighbor cell
				// is offset by the grid coordinates around the
				// current cell. In order to wrap the board, offet
				// by the dimension of the board and modulo to get
				// where it is.
				row = (cellRow + i + b.dim) % b.dim
				col = (cellCol + j + b.dim) % b.dim
			)

			// Don't count the current cell as a neighbor.
			if cellRow == row && cellCol == col {
				continue
			}

			if b.gen[row][col] == L {
				neighbors++
			}
		}
	}

	state := b.gen[cellRow][cellCol]
	switch neighbors {
	case 0, 1:
		return D
	case 2:
		if state == L {
			return L
		}
	case 3:
		return L
	}

	return D
}
