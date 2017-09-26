package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	allCandidates = genAllCandidates(NUM_COLS)
	os.Exit(t.Run())
}

func TestKnuthGuess(t *testing.T) {
	kG := knuthGuess(nil)
	firstRes := 1122
	if !reflect.DeepEqual(kG, firstRes) {
		t.Errorf("got: %v, want: %v", kG, firstRes)
	}
}

func TestScore(t *testing.T) {
	bulls, cows := score(1362, 3632)
	if bulls != 1 || cows != 2 {
		t.Errorf("expected 1 bull, 2 cows, got %d bulls and %d cows", bulls, cows)
	}
}

func TestGen(t *testing.T) {
	// This example comes from p3 of the Knuth mastermind paper.
	gotF := gen(3632)
	expectedF3632 := []feedback{
		{guess: 1122, bulls: 1, cows: 0},
		{guess: 1344, bulls: 0, cows: 1},
		{guess: 3526, bulls: 1, cows: 2},
		{guess: 1462, bulls: 1, cows: 1},
		{guess: 3632, bulls: 4, cows: 0},
	}
	if !reflect.DeepEqual(gotF, expectedF3632) {
		t.Errorf("got %v, expected %v", gotF, expectedF3632)
	}
}

func TestR(t *testing.T) {
	kk := knuth{}
	r(0, 0, 0, nil, 3632, &kk)
	expectedK3632 := knuth{
		move: 1122, bullsAndCows: map[int]knuth{5: {move: 1344, bullsAndCows: map[int]knuth{1: {move: 3526, bullsAndCows: map[int]knuth{7: {move: 1462, bullsAndCows: map[int]knuth{6: {move: 3632, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}}}
	if !reflect.DeepEqual(kk, expectedK3632) {
		t.Errorf("got %v, expected %v", kk, expectedK3632)
	}
}

func TestGen2(t *testing.T) {
	f := [][]feedback{}
	gen2(nil, &f)
	fmt.Println(f)
}

func TestSol(t *testing.T) {
	//fmt.Println(s)
	//fmt.Println(s.bullsAndCows[5].bullsAndCows[1].bullsAndCows[7].bullsAndCows[6])
}

func TestKnuthGen(t *testing.T) {
	expected := knuth{move: 1122, bullsAndCows: map[int]knuth{5: {move: 1344, bullsAndCows: map[int]knuth{1: {move: 3526, bullsAndCows: map[int]knuth{7: {move: 1462, bullsAndCows: map[int]knuth{6: {move: 3632, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}, 0: {move: 5525, bullsAndCows: map[int]knuth{2: {move: 6652, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}
	got := knuthSolutionGenerator([]int{3632, 6652})

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}

	got = knuthSolutionGenerator(
		allCandidates[:2],
	)
	expected = knuth{move: 1122, bullsAndCows: map[int]knuth{15: {move: 1223, bullsAndCows: map[int]knuth{6: {move: 1114, bullsAndCows: map[int]knuth{15: {move: 1112, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}, 10: {move: 1234, bullsAndCows: map[int]knuth{5: {move: 1315, bullsAndCows: map[int]knuth{10: {move: 1111, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func TestGenAndTime(t *testing.T) {
	start := time.Now()
	got := knuthSolutionGenerator(
		allCandidates,
	)
	fmt.Printf("%#v\n", got)
	//fmt.Println(got)
	fmt.Println(time.Since(start))
}

func TestMerge(t *testing.T) {
	// has 1111,1112,1113
	a := knuth{move: 1122, bullsAndCows: map[int]knuth{10: {move: 1234, bullsAndCows: map[int]knuth{5: {move: 1315, bullsAndCows: map[int]knuth{10: {move: 1111, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}, 6: {move: 2156, bullsAndCows: map[int]knuth{5: {move: 1113, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}, 15: {move: 1223, bullsAndCows: map[int]knuth{6: {move: 1114, bullsAndCows: map[int]knuth{15: {move: 1112, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}
	// has 1114
	b := knuth{move: 1122, bullsAndCows: map[int]knuth{10: {move: 1234, bullsAndCows: map[int]knuth{10: {move: 1536, bullsAndCows: map[int]knuth{5: {move: 1114, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}
	fmt.Println(a)
	merge(&a, b)
	expected := knuth{move: 1122, bullsAndCows: map[int]knuth{15: {move: 1223, bullsAndCows: map[int]knuth{6: {move: 1114, bullsAndCows: map[int]knuth{15: {move: 1112, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}, 10: {move: 1234, bullsAndCows: map[int]knuth{10: {move: 1536, bullsAndCows: map[int]knuth{5: {move: 1114, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}, 6: {move: 2156, bullsAndCows: map[int]knuth{5: {move: 1113, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}, 5: {move: 1315, bullsAndCows: map[int]knuth{10: {move: 1111, bullsAndCows: map[int]knuth{20: {move: 0, bullsAndCows: map[int]knuth(nil)}}}}}}}}}
	if !reflect.DeepEqual(a, expected) {
		t.Errorf("got %v\n\n, expected %v", a, expected)
	}
}
