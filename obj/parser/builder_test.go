package parser

import (
	"computer_graphics/obj/parser/types"
	"computer_graphics/obj/scanner"
	"testing"
)

// Testing a finite state machine table of an arbitrary elementParser.
func testParser(parser elementParser, want [][scanner.TokensCount]stateType, t *testing.T) {
	var (
		got     = parser.(*finiteStateMachine).matrix
		gotDim  = len(got)
		wantDim = len(want)
	)
	if gotDim != wantDim {
		t.Fatalf("Incorrect dimension of the matrix, got: %d, want: %d", gotDim, wantDim)
	}
	var correct = true
	for i := 0; i < gotDim; i++ {
		for j := 0; j < scanner.TokensCount; j++ {
			if got[i][j] != want[i][j] {
				t.Errorf(
					"Invalid matrix element (%d, %s), got: %d, want: %d",
					i,
					scanner.TokenType(j),
					got[i][j],
					want[i][j],
				)
				correct = false
			}
		}
	}
	if !correct {
		t.Log("got:  ", got)
		t.Log("want: ", want)
	}
}

// Testing the vertex elementParser.
func TestBuildParser_vertex(t *testing.T) {
	var (
		parser = buildParser(Vertex, types.NewVertex())
		want   = [][scanner.TokensCount]stateType{
			{1, 1, 1, 1, 2, 1, 1, 1, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 3, 3, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 4, 1, 1, 1, 1},
			{1, 5, 5, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 6, 1, 1, 1, 1},
			{1, 7, 7, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 8, 0, 0, 1, 1},
			{1, 9, 9, 1, 1, 0, 0, 1, 1},
			{1, 1, 1, 1, 10, 0, 0, 1, 1},
			{1, 1, 1, 1, 1, 0, 0, 1, 1},
		}
	)
	testParser(parser, want, t)
}

// Testing the face elementParser.
func TestBuildParser_face(t *testing.T) {
	var (
		parser = buildParser(Face, types.NewFace())
		want   = [][scanner.TokensCount]stateType{
			{1, 1, 1, 1, 2, 1, 1, 1, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 3, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 4, 55, 1, 1, 1, 1},
			{1, 5, 1, 38, 1, 1, 1, 1, 1},
			{1, 1, 1, 6, 26, 1, 1, 1, 1},
			{1, 7, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 8, 1, 1, 1, 1},
			{1, 9, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 10, 1, 1, 1, 1, 1},
			{1, 11, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 12, 1, 1, 1, 1, 1},
			{1, 13, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 14, 1, 1, 1, 1},
			{1, 15, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 16, 1, 1, 1, 1, 1},
			{1, 17, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 18, 1, 1, 1, 1, 1},
			{1, 19, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 20, 0, 0, 1, 1},
			{1, 21, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 22, 1, 1, 1, 1, 1},
			{1, 23, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 24, 1, 1, 1, 1, 1},
			{1, 25, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 20, 0, 0, 1, 1},
			{1, 27, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 28, 1, 1, 1, 1, 1},
			{1, 29, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 30, 1, 1, 1, 1},
			{1, 31, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 32, 1, 1, 1, 1, 1},
			{1, 33, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 34, 0, 0, 1, 1},
			{1, 35, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 36, 1, 1, 1, 1, 1},
			{1, 37, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 34, 0, 0, 1, 1},
			{1, 39, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 40, 1, 1, 1, 1},
			{1, 41, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 42, 1, 1, 1, 1, 1},
			{1, 1, 1, 43, 1, 1, 1, 1, 1},
			{1, 44, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 45, 1, 1, 1, 1},
			{1, 46, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 47, 1, 1, 1, 1, 1},
			{1, 1, 1, 48, 1, 1, 1, 1, 1},
			{1, 49, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 50, 0, 0, 1, 1},
			{1, 51, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 52, 1, 1, 1, 1, 1},
			{1, 1, 1, 53, 1, 1, 1, 1, 1},
			{1, 54, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 50, 0, 0, 1, 1},
			{1, 56, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 57, 1, 1, 1, 1},
			{1, 58, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 59, 0, 0, 1, 1},
			{1, 60, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 59, 0, 0, 1, 1},
		}
	)
	testParser(parser, want, t)
}
