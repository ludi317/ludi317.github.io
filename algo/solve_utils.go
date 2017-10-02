package main

import "sort"

func merge(m1 knuth, m2 knuth) {
	for k2, v2 := range m2.next {
		if v1, ok := m1.next[k2]; ok {
			merge(v1, v2)
		} else {
			(m1).next[k2] = v2
		}
	}
}
func genAllFeedback() []int {
	res := make([]int, (numCols+1)*(numCols+2)/2-1)
	c := 0
	for bulls := numCols; bulls >= 0; bulls-- {
		for cows := 0; cows <= numCols-bulls; cows++ {
			if bulls == numCols-1 && cows == 1 {
				continue
			}
			res[c] = hash(bulls, cows)
			c++
		}
	}
	return res
}

func score(guess int, solution int) int {
	bulls, cows := 0, 0
	colors := make([]int, numColors+1)
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

func hash(bulls, cows int) int {
	return (numCols+1)*bulls + cows
}

func genAllCandidates() []int {
	res := genAllCandidatesHelper(numCols)
	sort.Ints(res)
	return res
}

func genAllCandidatesHelper(numCols int) []int {
	if numCols == 1 {
		cands := make([]int, numColors)
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
