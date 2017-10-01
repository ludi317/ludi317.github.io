package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func main() {
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
		panic(fmt.Sprintf("wrong solution\ngot:\n%v\n\nwant:\n%v", got, kSol))
	}
	fmt.Println()
}
