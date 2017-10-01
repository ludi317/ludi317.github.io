package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const numCols = 4
const numColors = 6
const size = 8

func main() {
	start := time.Now()
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
	fmt.Println(time.Since(start), "colors:", numColors, "columns:", numCols, "batchsize:", size)
	if numCols == 4 && numColors == 6 && !reflect.DeepEqual(got, kSol) {
		panic(fmt.Sprintf("wrong solution\ngot:\n%v\n\nwant:\n%v", got, kSol))
	}
	fmt.Println()
}
