package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	numCols   = 4
	numColors = 6
)

func main() {
	start := time.Now()
	solutionTrie := knuthSolutionGeneratorIter()
	t := time.Since(start)

	f, err := os.Create(fmt.Sprintf("colors_%d_cols_%d.txt", numColors, numCols))
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	fmt.Println(t, "colors:", numColors, "columns:", numCols, "maxMoves:", maxMoves)

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, t, "colors:", numColors, "columns:", numCols, "maxMoves:", maxMoves)
	s := fmt.Sprintf("%#v\n", solutionTrie)
	s = strings.Replace(s, "main.", "", -1)
	fmt.Fprintf(w, s)

	if numCols == 4 && numColors == 6 && !reflect.DeepEqual(solutionTrie, kSol) {
		panic(fmt.Sprintf("wrong solution\ngot:\n%v\n\nwant:\n%v", solutionTrie, kSol))
	}
	w.Flush()
}
