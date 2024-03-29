package main

var (
	allCodes    = genAllCodes()
	allFeedback = genAllFeedback()
	maxMoves    int
	allBulls    = hash(numCols, 0)
)

type feedback struct {
	guess int
	bc    int
	skip  bool
}

type knuth struct {
	move int
	next map[int]knuth
}

func generateKnuthBranchIter(solution int, total knuth) knuth {
	var feedbacks []feedback
	bc := 0
	guess := 0
	invalids := make([]bool, len(allCodes))
	for bc != allBulls {
		next, ok := total.next[bc]
		if ok {
			guess = next.move
			total = next
		} else {
			guess = knuthGuess(feedbacks, &invalids)
			total = knuth{}
		}
		bc = score(guess, solution)
		feedbacks = append(feedbacks, feedback{guess: guess, bc: bc})
	}
	kk := knuth{}
	if len(feedbacks) > maxMoves {
		maxMoves = len(feedbacks)
	}
	for i := len(feedbacks) - 1; i >= 0; i-- {
		kk.next = map[int]knuth{feedbacks[i].bc: kk}
		kk.move = feedbacks[i].guess
	}
	return kk
}

// knuthSolutionGeneratorIter generates a trie of knuth structs that records the move to make for all possible solutions.
// Each node stores the move to make, and a map. The keys of the map span the range of possible feedback and values are
// downstream nodes. Output is stored as a variable in solutions.go.
func knuthSolutionGeneratorIter() knuth {

	total := knuth{next: map[int]knuth{}}
	batchSize := getBatchSize()
	shuffledCodes := shuffle()
	ch := make(chan knuth, batchSize)
	ch2 := make(chan bool, batchSize)
	batches := len(allCodes) / batchSize
	for i := 0; i < batches; i++ {
		for j := i * batchSize; j < (i+1)*batchSize; j++ {
			go func(val int) {
				// Concurrent reads on the total trie.
				kk := generateKnuthBranchIter(val, total)
				kk = knuth{next: map[int]knuth{0: kk}}
				ch <- kk
				ch2 <- true
			}(shuffledCodes[j])
		}
		for i := 0; i < batchSize; i++ {
			<-ch2
		}
		// All goroutines have finished by now.
		// At this point it is safe to write to the total trie.
		for i := 0; i < batchSize; i++ {
			kk := <-ch
			merge(total, kk)
		}
	}
	return total.next[0]
}

// isValid indicates if the candidate is a possible valid solution given the feedback.
func isValid(c int, fs []feedback) bool {
	for _, f := range fs {
		if !f.skip {
			if score(f.guess, c) != f.bc {
				return false
			}
		}
	}
	return true
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
func knuthGuess(feedbacks []feedback, invalids *[]bool) int {
	// Calculate max number of remaining possibilities over all feedback.
	scores := maxSolutionSpaceSize(&feedbacks, invalids)

	// Initialize minScore to its highest possible value.
	minScore := len(allCodes)

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
		if !(*invalids)[pos] {
			return allCodes[pos]
		}
	}
	// out of range index panic here means there were no possible solutions given feedback
	return allCodes[candMinScoresPos[0]]
}

func maxSolutionSpaceSize(fs *[]feedback, invalids *[]bool) []int {
	scores := make([]int, len(allCodes))
	for i, hypoGuess := range allCodes {
		possibleScores := make([]int, (numCols+1)*(numCols+1))
		for idx, hypoSoln := range allCodes {
			if (*invalids)[idx] {
				continue
			}
			if !isValid(hypoSoln, *fs) {
				(*invalids)[idx] = true
				continue
			}
			possibleScores[score(hypoGuess, hypoSoln)]++
		}
		for i := range *fs {
			if !(*fs)[i].skip {
				(*fs)[i].skip = true
			}
		}
		maxSS := 0
		for _, s := range possibleScores {
			if s > maxSS {
				maxSS = s
			}
		}
		scores[i] = maxSS
	}
	return scores
}

// genKnuthBranchRec is a recursive implementation that creates a single branch of the knuth trie.
func genKnuthBranchRec(bc int, guess int, fs []feedback, solution int, kk *knuth, total knuth, invalids *[]bool) {
	if bc == hash(numCols, 0) {
		return
	}
	next, ok := total.next[bc]
	if ok {
		guess = next.move
		total = next
	} else {
		guess = knuthGuess(fs, invalids)
		total = knuth{}
	}
	bc = score(guess, solution)
	fs = append(fs, feedback{guess: guess, bc: bc})
	genKnuthBranchRec(bc, guess, fs, solution, kk, total, invalids)
	*kk = knuth{move: guess, next: map[int]knuth{bc: *kk}}
}
