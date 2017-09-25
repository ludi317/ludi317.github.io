package main

func score(guess []int, solution []int) (int, int) {
	bulls, cows := 0, 0
	colors := [NUM_COLORS + 1]int{}
	for i, g := range guess {
		if solution[i] == g {
			bulls++
		} else {
			if colors[g] < 0 {
				cows++
			}
			if colors[solution[i]] > 0 {
				cows++
			}
			colors[g]++
			colors[solution[i]]--
		}
	}
	return bulls, cows
}

func allFeedback() []feedback {
	res := []feedback{}
	for bulls := 0; bulls <= NUM_COLS; bulls++ {
		for cows := 0; cows <= NUM_COLS-bulls; cows++ {
			if bulls == 3 && cows == 1 {
				continue
			}
			res = append(res, feedback{bulls: bulls, cows: cows})
		}
	}
	return res
}

func calcNum(hf feedback, ch chan int) {
	numPoss := numPossibilities(hf)
	ch <- numPoss
}

type s struct {
	i       int
	maxPoss int
}

// Calculate # of remaining possibilities for every type of feedback. Score is the max of that.
func calcMax(c candidate, ind int, ch2 chan s) {

	ch := make(chan int)

	aF := allFeedback()
	for _, f := range aF {

		hf := feedback{
			guess: c.code,
			bulls: f.bulls,
			cows:  f.cows,
		}
		go calcNum(hf, ch)
	}
	maxPoss := <-ch
	for i := 0; i < len(aF)-1; i++ {
		numPoss := <-ch
		if numPoss > maxPoss {
			maxPoss = numPoss
		}
	}
	ch2 <- s{maxPoss: maxPoss, i: ind}
}

func knuthGuess() []int {
	if activeRow == 0 {
		return []int{1, 1, 2, 2}
	}
	ch2 := make(chan s)

	for ind, c := range allCandidates {
		go calcMax(c, ind, ch2)
	}
	for i := 0; i < len(allCandidates); i++ {
		s := <-ch2
		allCandidates[s.i].score = s.maxPoss
	}

	// initialize minScore to its highest possible value: NUM_COLORS^NUM_COLS
	minScore := 1
	for i := 0; i < NUM_COLS; i++ {
		minScore *= NUM_COLORS
	}

	candMinScores := []candidate{}
	for _, c := range allCandidates {
		if c.score < minScore {
			minScore = c.score
			candMinScores = []candidate{c}
		} else if c.score == minScore {
			candMinScores = append(candMinScores, c)
		}
	}

	for _, c := range candMinScores {
		if isValid(c, feedbacks) {
			return c.code
		}
	}
	return candMinScores[0].code
}

func numPossibilities(hf feedback) int {
	f := make([]feedback, len(feedbacks)+1)
	copy(f, feedbacks)
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
func isValid(c candidate, fs []feedback) bool {
	for _, f := range fs {
		if bulls, cows := score(f.guess, c.code); bulls != f.bulls || cows != f.cows {
			return false
		}
	}
	return true
}

func genAllCandidates(numCols int) []candidate {
	if numCols == 1 {
		cands := make([]candidate, NUM_COLORS)
		for i := range cands {
			cands[i].code = []int{i + 1}
		}
		return cands
	}
	return product(genAllCandidates(1), genAllCandidates(numCols-1))
}

func product(a []candidate, b []candidate) []candidate {
	res := make([]candidate, len(a)*len(b))
	c := 0
	for _, ae := range a {
		for _, be := range b {
			res[c].code = concat(ae.code, be.code...)
			c++
		}
	}
	return res
}

func concat(a []int, b ...int) []int {
	c := make([]int, len(a))
	copy(c, a)
	c = append(c, b...)
	return c
}
