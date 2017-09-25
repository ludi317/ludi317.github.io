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

// genGuess returns the first guess from a sorted list that could be the solution given the feedback.
func genGuess() []int {
OuterLoop:
	for i, c := range candidates {
		if c == nil {
			continue
		}
		for _, f := range feedbacks {
			if bulls, cows := score(f.guess, c); bulls != f.bulls || cows != f.cows {
				candidates[i] = nil
				continue OuterLoop
			}
		}
		return c
	}
	return nil
}

func allCandidates(numCols int) [][]int {
	if numCols == 1 {
		a := make([][]int, NUM_COLORS)
		for i := range a {
			a[i] = []int{i + 1}
		}
		return a
	}
	return product(allCandidates(1), allCandidates(numCols-1))
}

func product(a [][]int, b [][]int) [][]int {
	res := make([][]int, len(a)*len(b))
	c := 0
	for _, ae := range a {
		for _, be := range b {
			res[c] = concat(ae, be...)
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
