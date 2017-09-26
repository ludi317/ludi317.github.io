package main

import (
	"fmt"
	"sort"
)

type knuth struct {
	move         int
	bullsAndCows map[int]knuth
}

func score(guess int, solution int) (int, int) {
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
	return bulls, cows
}

func allFeedback() []feedback {
	res := []feedback{}
	for bulls := NUM_COLS; bulls >= 0; bulls-- {
		for cows := 0; cows <= NUM_COLS-bulls; cows++ {
			if bulls == 3 && cows == 1 {
				continue
			}
			res = append(res, feedback{bulls: bulls, cows: cows})
		}
	}
	return res
}


// fortran array indexing
func hash(bulls, cows int) int {
	return (NUM_COLS+1)*bulls + cows
}

// knuthGuess is implementation of Knuth algo that guarantees the solution in <= 5 moves.
// See: http://www.cs.uni.edu/~wallingf/teaching/cs3530/resources/knuth-mastermind.pdf
func knuthGuess(feedbacks []feedback) int {
	if len(feedbacks) == 0 {
		return 1122
	}
	valid := false
	scores := make([]int, len(allCandidates))
	for i, c := range allCandidates {
		maxPoss := 0
		// Calculate # of remaining possibilities for every type of feedback. Score is the max of that.
		for _, af := range allFeedback() {
			hf := feedback{
				guess: c,
				bulls: af.bulls,
				cows:  af.cows,
			}
			numPoss := numPossibilities(hf, feedbacks)
			if numPoss > maxPoss {
				maxPoss = numPoss
				valid = true
			}
		}
		scores[i] = maxPoss
	}

	if !valid {
		return -1
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
	return -1
}

// Forward pass
func gen(solution int) []feedback {
	bulls, cows := 0, 0
	var feedbacks []feedback

	for bulls != NUM_COLS {
		guess := knuthGuess(feedbacks)
		bulls, cows = score(guess, solution)
		feedbacks = append(feedbacks, feedback{guess: guess, bulls: bulls, cows: cows})
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
	if len(feed) > 0 && feed[len(feed)-1].bulls == NUM_COLS {
		*keep = append(*keep, feed)
		fmt.Println(*keep)
		return true
	}
	for _, f := range allFeedback() {
		copyFeed := make([]feedback, len(feed)+1)
		copy(copyFeed, feed)
		copyFeed[len(copyFeed)-1] = f

		copyFeed[len(copyFeed)-1].guess = guess
		fmt.Println(copyFeed)
		if gen2(copyFeed, keep) {
			break
		}
	}
	return true
}

// knuthSolutionGenerator generates a trie of knuth structs that records the move to make for all possible solutions.
// Each node stores the move to make, and a map. The keys of the map span the range of possible feedback and
// values are downstream nodes. Output is stored as a variable in solutions.go. It takes ~13 mins to run on my laptop.
func knuthSolutionGenerator(cs []int) knuth {

	total := knuth{move: 1122, bullsAndCows: map[int]knuth{}}
	ch := make(chan knuth)
	s := 16
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
				r(0, 0, 0, nil, val, &kk)
				ch <- kk
				fmt.Println(val)
			}(cs[j])

		}
		for i := 0; i < s; i++ {
			kk := <-ch
			merge(&total, kk)
		}
	}
	return total
}

func merge(m1 *knuth, m2 knuth) {
	for k2, v2 := range m2.bullsAndCows {
		for k1, v1 := range m1.bullsAndCows {
			if k1 == k2 {
				m1 = &v1
				m2 = v2
				merge(m1, m2)
				return
			}
		}
		(*m1).bullsAndCows[k2] = v2
		return
	}
}

// r is a recursive implementation that creates a single branch of the knuth trie.
func r(cows, bulls int, guess int, fs []feedback, solution int, kk *knuth) {
	if bulls == NUM_COLS {
		return
	}
	guess = knuthGuess(fs)
	bulls, cows = score(guess, solution)
	fs = append(fs, feedback{guess: guess, bulls: bulls, cows: cows})
	r(cows, bulls, guess, fs, solution, kk)
	*kk = knuth{move: guess, bullsAndCows: map[int]knuth{hash(bulls, cows): *kk}}
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
		if bulls, cows := score(f.guess, c); bulls != f.bulls || cows != f.cows {
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
