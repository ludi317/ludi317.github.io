package main

import (
	"fmt"
	"sort"
)

type feedback struct {
	guess int
	bc    int
}

type knuth struct {
	move int
	next map[int]knuth
}

func score(guess int, solution int) int {
	bulls, cows := 0, 0
	colors := [NUM_COLORS + 1]int{}
	for ; guess > 0; solution, guess = solution/10, guess/10 {
		sCol := solution % 10
		gCol := guess % 10
		if sCol == gCol {
			bulls++
		} else {
			if colors[gCol] < 0 {
				cows++
			}
			if colors[sCol] > 0 {
				cows++
			}
			colors[gCol]++
			colors[sCol]--
		}
	}
	return hash(bulls, cows)
}

func allFeedback() []int {
	res := make([]int, (NUM_COLS+1)*(NUM_COLS+2)/2-1)
	c := 0
	for bulls := NUM_COLS; bulls >= 0; bulls-- {
		for cows := 0; cows <= NUM_COLS-bulls; cows++ {
			if bulls == NUM_COLS-1 && cows == 1 {
				continue
			}
			res[c] = hash(bulls, cows)
			c++
		}
	}
	return res
}

// fortran array indexing
func hash(bulls, cows int) int {
	return (NUM_COLS+1)*bulls + cows
}

// knuthGuess is implementation of Knuth algorithm that guarantees the solution in <= 5 moves. See:
// http://www.cs.uni.edu/~wallingf/teaching/cs3530/resources/knuth-mastermind.pdf
//
// It plays the code that gives the most information, where information is defined as how much the solution space is
// reduced. At every turn, it sweeps through all hypothetical feedback for all codes, counting the number of possible
// solutions remaining if that feedback were given. This number depends on all the feedback already present from prior
// moves. Each code is scored with the max solution space size of all hypothetical feedback. The code with the smallest
// score is played. To break ties, codes that are themselves possible solutions are preferred, followed by numerical
// ordering. More concisely, the code chosen has the min of the max of the possible remaining solutions.
func knuthGuess(feedbacks []feedback) int {
	scores := make([]int, len(allCandidates))
	for i, c := range allCandidates {
		maxPoss := 0
		// Calculate # of remaining possibilities for every type of feedback. Score is the max of that.
		for _, af := range allFeedback() {
			hf := feedback{
				guess: c,
				bc:    af,
			}
			numPoss := numPossibilities(hf, feedbacks)
			if numPoss > maxPoss {
				maxPoss = numPoss
			}
		}
		scores[i] = maxPoss
	}

	// Initialize minScore to its highest possible value: NUM_COLORS^NUM_COLS
	minScore := 1
	for i := 0; i < NUM_COLS; i++ {
		minScore *= NUM_COLORS
	}

	// Choose the candidate that minimizes the (max) remaining possibilities.
	candMinScoresPos := []int{}
	for i, s := range scores {
		if s < minScore {
			minScore = s
			candMinScoresPos = []int{i}
		} else if s == minScore {
			candMinScoresPos = append(candMinScoresPos, i)
		}
	}

	// Prefer candidates that could be the solution.
	for _, pos := range candMinScoresPos {
		if isValid(allCandidates[pos], feedbacks) {
			return allCandidates[pos]
		}
	}
	if len(candMinScoresPos) > 0 {
		return allCandidates[candMinScoresPos[0]]
	}
	panic(fmt.Sprintf("no possible solutions given feedback %v", feedbacks))
}

// Forward pass
func gen(solution int) []feedback {
	var feedbacks []feedback
	bc := -1
	for bc != hash(NUM_COLS, 0) {
		guess := knuthGuess(feedbacks)
		bc = score(guess, solution)
		feedbacks = append(feedbacks, feedback{guess: guess, bc: bc})
		fmt.Printf("%#v\n", feedbacks)
	}
	return feedbacks
}

func gen2(feed []feedback, keep *[][]feedback) bool {
	guess := knuthGuess(feed)
	if guess == -1 {
		return false
	}
	fmt.Println(guess)
	if len(feed) > 0 && feed[len(feed)-1].bc == hash(NUM_COLS, 0) {
		*keep = append(*keep, feed)
		fmt.Println(*keep)
		return true
	}
	for _, f := range allFeedback() {
		copyFeed := make([]feedback, len(feed)+1)
		copy(copyFeed, feed)
		copyFeed[len(copyFeed)-1] = feedback{bc: f}

		copyFeed[len(copyFeed)-1].guess = guess
		fmt.Println(copyFeed)
		if gen2(copyFeed, keep) {
			break
		}
	}
	return true
}

// knuthSolutionGenerator generates a trie of knuth structs that records the move to make for all possible solutions.
// Each node stores the move to make, and a map. The keys of the map span the range of possible feedback and values are
// downstream nodes. Output is stored as a variable in solutions.go.
func knuthSolutionGenerator(cs []int, s int) knuth {

	total := knuth{next: map[int]knuth{}}
	ch := make(chan knuth)
	if s > len(cs) {
		s = len(cs)
	}
	if len(cs)%s != 0 {
		panic("bad batch size")
	}
	batches := len(cs) / s
	for i := 0; i < batches; i++ {
		for j := i * s; j < (i+1)*s; j++ {
			go func(val int) {
				kk := knuth{}
				r(0, 0, nil, val, &kk, total)
				kk = knuth{next: map[int]knuth{0: kk}}
				ch <- kk
				fmt.Println(val)
			}(cs[j])

		}
		for i := 0; i < s; i++ {
			kk := <-ch
			merge(&total, kk)
		}
	}
	return total.next[0]
}

func merge(m1 *knuth, m2 knuth) {
	for k2, v2 := range m2.next {
		if v1, ok := m1.next[k2]; ok {
			m1 = &v1
			m2 = v2
			merge(m1, m2)
		} else {
			(*m1).next[k2] = v2
		}
	}
}

// r is a recursive implementation that creates a single branch of the knuth trie.
func r(bc int, guess int, fs []feedback, solution int, kk *knuth, total knuth) {
	if bc == hash(NUM_COLS, 0) {
		return
	}
	if next, ok := total.next[bc]; ok && guess == total.move {
		guess = next.move
		total = next
	} else {
		guess = knuthGuess(fs)
		total = knuth{}
	}
	bc = score(guess, solution)
	fs = append(fs, feedback{guess: guess, bc: bc})
	r(bc, guess, fs, solution, kk, total)
	*kk = knuth{move: guess, next: map[int]knuth{bc: *kk}}
}

func numPossibilities(hf feedback, fs []feedback) int {
	f := make([]feedback, len(fs)+1)
	copy(f, fs)
	f[len(f)-1] = hf
	n := 0
	for _, c := range allCandidates {
		if isValid(c, f) {
			n++
		}
	}
	return n
}

// isValid indicates if the candidate is a possible valid solution given the feedback.
func isValid(c int, fs []feedback) bool {
	for _, f := range fs {
		if score(f.guess, c) != f.bc {
			return false
		}
	}
	return true
}

func genAllCandidates(numCols int) []int {
	res := genAllCandidatesHelper(numCols)
	sort.Ints(res)
	return res
}

func genAllCandidatesHelper(numCols int) []int {
	if numCols == 1 {
		cands := make([]int, NUM_COLORS)
		for i := range cands {
			cands[i] = i + 1
		}
		return cands
	}
	return product(genAllCandidatesHelper(1), genAllCandidatesHelper(numCols-1))
}

func product(a []int, b []int) []int {
	res := make([]int, len(a)*len(b))
	c := 0
	for _, ae := range a {
		for _, be := range b {
			res[c] = ae + be*10
			c++
		}
	}
	return res
}
