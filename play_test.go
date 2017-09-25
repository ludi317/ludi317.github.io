package main

import (
	"reflect"
	"testing"
)

func TestKnuthGuess(t *testing.T) {
	allCandidates = genAllCandidates(NUM_COLS)
	kG := knuthGuess()
	firstRes := []int{1, 1, 2, 2}
	if !reflect.DeepEqual(kG, firstRes) {
		t.Errorf("got: %v, want: %v", kG, firstRes)
	}

}
