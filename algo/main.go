package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
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
	f, err := os.Create(fmt.Sprintf("colors_%d_cols_%d.txt", numColors, numCols))
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Println(time.Since(start), "colors:", numColors, "columns:", numCols, "maxMoves:", maxMoves, "batchsize:", size)
	fmt.Fprintln(w, time.Since(start), "colors:", numColors, "columns:", numCols, "maxMoves:", maxMoves, "batchsize:", size)
	s := fmt.Sprintf("%#v\n", got)
	s = strings.Replace(s, "main.", "", -1)
	fmt.Fprintf(w, s)
	if numCols == 4 && numColors == 6 && !reflect.DeepEqual(got, kSol) {
		panic(fmt.Sprintf("wrong solution\ngot:\n%v\n\nwant:\n%v", got, kSol))
	}
	w.Flush()
	fmt.Println()
}
