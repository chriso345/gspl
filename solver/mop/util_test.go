package mop

import "testing"

func TestAlmostEqual(t *testing.T) {
	a := []float64{1.0, 2.0}
	b := []float64{1.0 + 1e-10, 2.0 - 1e-10}

	if !almostEqual(a, b) {
		t.Fatalf("expected vectors to be almost equal")
	}
}

func TestUniqueSolutions(t *testing.T) {
	sols := []*MopSolution{
		{ObjectiveValues: []float64{1, 2}},
		{ObjectiveValues: []float64{1, 2}},
		{ObjectiveValues: []float64{2, 1}},
	}

	uniq := uniqueSolutions(sols)

	if len(uniq) != 2 {
		t.Fatalf("expected 2 unique solutions, got %d", len(uniq))
	}

	if uniq[0].ObjectiveValues[0] > uniq[1].ObjectiveValues[0] {
		t.Fatalf("solutions not sorted by first objective")
	}
}
