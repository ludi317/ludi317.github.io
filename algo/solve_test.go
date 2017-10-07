package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestKnuthGuess(t *testing.T) {
	invalids := make([]bool, len(allCodes))
	kG := knuthGuess(nil, &invalids)
	firstRes := 1122
	if !reflect.DeepEqual(kG, firstRes) {
		t.Errorf("got: %v, want: %v", kG, firstRes)
	}
}

func TestScore(t *testing.T) {
	bc := score(1362, 3632)
	if bc != hash(1, 2) {
		t.Errorf("expected: %d, got: %d", hash(1, 2), bc)
	}
}

func TestKnuthBranchIter(t *testing.T) {
	// This example comes from p3 of the Knuth mastermind paper.
	gotF := generateKnuthBranchIter(3632, knuth{})
	expectedF3632 := knuth{move: 1122, next: map[int]knuth{5: {move: 1344, next: map[int]knuth{1: {move: 3526, next: map[int]knuth{7: {move: 1462, next: map[int]knuth{6: {move: 3632, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}}}
	if !reflect.DeepEqual(gotF, expectedF3632) {
		t.Errorf("got %#v, expected %#v", gotF, expectedF3632)
	}
}

func TestKnuthBranchRec(t *testing.T) {
	expectedK3632 := knuth{
		move: 1122, next: map[int]knuth{5: {move: 1344, next: map[int]knuth{1: {move: 3526, next: map[int]knuth{7: {move: 1462, next: map[int]knuth{6: {move: 3632, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}}}
	//total := knuth{0, map[int]knuth{0: expectedK3632}}
	total := expectedK3632
	kk := knuth{}
	valids := make([]bool, len(allCodes))
	genKnuthBranchRec(0, 0, nil, 3632, &kk, total, &valids)
	if !reflect.DeepEqual(kk, expectedK3632) {
		t.Errorf("got %v, expected %v", kk, expectedK3632)
	}
}

func TestGenAndTime(t *testing.T) {
	start := time.Now()
	got := knuthSolutionGeneratorIter()
	fmt.Println(time.Since(start))
	if numCols == 4 && numColors == 6 && !reflect.DeepEqual(got, kSol) {
		t.Errorf("wrong solution\ngot:\n%v\n\nwant:\n%v", got, kSol)
	}
}

func TestMerge(t *testing.T) {
	// has 1111,1112,1113
	a := knuth{move: 1122, next: map[int]knuth{10: {move: 1234, next: map[int]knuth{5: {move: 1315, next: map[int]knuth{10: {move: 1111, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}, 6: {move: 2156, next: map[int]knuth{5: {move: 1113, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}, 15: {move: 1223, next: map[int]knuth{6: {move: 1114, next: map[int]knuth{15: {move: 1112, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}
	// has 1114
	b := knuth{move: 1122, next: map[int]knuth{10: {move: 1234, next: map[int]knuth{10: {move: 1536, next: map[int]knuth{5: {move: 1114, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}
	fmt.Println(a)
	merge(a, b)
	expected := knuth{move: 1122, next: map[int]knuth{15: {move: 1223, next: map[int]knuth{6: {move: 1114, next: map[int]knuth{15: {move: 1112, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}, 10: {move: 1234, next: map[int]knuth{10: {move: 1536, next: map[int]knuth{5: {move: 1114, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}, 6: {move: 2156, next: map[int]knuth{5: {move: 1113, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}, 5: {move: 1315, next: map[int]knuth{10: {move: 1111, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}
	if !reflect.DeepEqual(a, expected) {
		t.Errorf("got %v\n\n, expected %v", a, expected)
	}
}
