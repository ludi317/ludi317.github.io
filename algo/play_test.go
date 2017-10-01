package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestKnuthGuess(t *testing.T) {
	kG := knuthGuess(nil)
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

func TestGen(t *testing.T) {
	// This example comes from p3 of the Knuth mastermind paper.
	gotF := generateKnuthBranchIterNoCache(3632)
	expectedF3632 := []feedback{
		{guess: 1122, bc: hash(1, 0)},
		{guess: 1344, bc: hash(0, 1)},
		{guess: 3526, bc: hash(1, 2)},
		{guess: 1462, bc: hash(1, 1)},
		{guess: 3632, bc: hash(4, 0)},
	}
	if !reflect.DeepEqual(gotF, expectedF3632) {
		t.Errorf("got %v, expected %v", gotF, expectedF3632)
	}
}

func TestR(t *testing.T) {
	expectedK3632 := knuth{
		move: 1122, next: map[int]knuth{5: {move: 1344, next: map[int]knuth{1: {move: 3526, next: map[int]knuth{7: {move: 1462, next: map[int]knuth{6: {move: 3632, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}}}
	//total := knuth{0, map[int]knuth{0: expectedK3632}}
	total := expectedK3632
	kk := knuth{}
	genKnuthBranchRec(0, 0, nil, 3632, &kk, total)
	if !reflect.DeepEqual(kk, expectedK3632) {
		t.Errorf("got %v, expected %v", kk, expectedK3632)
	}
	fmt.Println(total)
}

func TestKnuthGen(t *testing.T) {
	expected := knuth{move: 1122, next: map[int]knuth{5: {move: 1344, next: map[int]knuth{1: {move: 3526, next: map[int]knuth{7: {move: 1462, next: map[int]knuth{6: {move: 3632, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}, 0: {move: 5525, next: map[int]knuth{2: {move: 6652, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}
	got := knuthSolutionGeneratorIter([]int{3632, 6652}, 1)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got \n%v, expected \n%v", got, expected)
	}

	got = knuthSolutionGeneratorIter(
		allCandidates[:2],
		2,
	)
	fmt.Println(got)
	expected = knuth{move: 1122, next: map[int]knuth{15: {move: 1223, next: map[int]knuth{6: {move: 1114, next: map[int]knuth{15: {move: 1112, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}, 10: {move: 1234, next: map[int]knuth{5: {move: 1315, next: map[int]knuth{10: {move: 1111, next: map[int]knuth{20: {move: 0, next: map[int]knuth(nil)}}}}}}}}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got \n%v, expected \n%v", got, expected)
	}
}

func TestGenAndTime(t *testing.T) {
	start := time.Now()
	size := 8
	shuffle := true
	cs := allCandidates
	if shuffle {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		shuffled := genAllCandidates()
		for i := len(shuffled) - 1; i > 0; i-- {
			rando := r.Intn(i + 1)
			shuffled[rando], shuffled[i] = shuffled[i], shuffled[rando]
		}
		cs = shuffled
	}
	got := knuthSolutionGeneratorIter(
		cs,
		size,
	)
	fmt.Println(got)
	fmt.Println(time.Since(start), "batchsize:", size, "with shuffling:", shuffle)
	if !reflect.DeepEqual(got, kSol) {
		t.Errorf("wrong solution\ngot:\n%v\n\nwant:\n%v", got, kSol)
	}
	fmt.Println()
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
