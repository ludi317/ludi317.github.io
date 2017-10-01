package main

type knuth struct {
	move int
	next map[int]knuth
}

func score(guess int, solution int) int {
	bulls, cows := 0, 0
	colors := [numColors + 1]int{}
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
func reverseHash(hash int) (bulls, cows int) {
	return hash / (numCols + 1), hash % (numCols + 1)
}
