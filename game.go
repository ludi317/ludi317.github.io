package main

import (
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/html/atom"
)

func generateSolution() int {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := 0
	for i := 0; i < numCols; i++ {
		s += rand1.Intn(numColors) + 1
		if i != numCols-1 {
			s *= 10
		}
	}
	return s
}

func pickColor(s int) {
	selectedColor = s
	document.GetElementByID(gameTableID).SetAttribute(atom.Style.String(), `cursor: url("`+imageDir+`color_`+
		strconv.Itoa(selectedColor)+`.gif"), auto; border: 3px black solid;`)
	document.GetElementByID("colorPickerID").SetAttribute(atom.Style.String(), `cursor: url("`+imageDir+`color_`+
		strconv.Itoa(selectedColor)+`.gif"), auto;`)
}

func placeColor(row, col int) {
	if selectedColor != 0 && row == activeRow {
		document.GetElementByID(strconv.Itoa(row)+"-"+strconv.Itoa(col)).SetAttribute(atom.Src.String(),
			imageDir+"color_"+strconv.Itoa(selectedColor)+".gif")
		guess = updateGuess(col)
		grade()
	}
}

func updateGuess(col int) int {
	m := 1
	copyGuess := guess
	for i := 0; i < numCols-col-1; i++ {
		m *= 10
		copyGuess /= 10
	}
	curGuess := copyGuess % 10
	return guess + m*(selectedColor-curGuess)
}

func solve() {
	go func() {
		for activeRow != -1 {
			m := kSol.move
			for col := numCols - 1; col >= 0; col-- {
				selectedColor = m % 10
				m /= 10
				placeColor(activeRow, col)
			}
			kSol = kSol.next[feedbackHash]
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func reload() {
	document.Location().Call("reload")
}

func hasZeros(guess int) bool {
	for i := 0; i < numCols; i++ {
		if guess%10 == 0 {
			return true
		}
		guess /= 10
	}
	return false
}

func grade() {
	if hasZeros(guess) {
		return
	}

	feedbackHash = score(guess, solution)
	bulls, cows := reverseHash(feedbackHash)
	pegHoles := document.GetElementsByClassName("graderRow" + strconv.Itoa(activeRow))
	i := 0
	for ; i < cows; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_6.gif")
	}
	for ; i < cows+bulls; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_5.gif")
	}

	if bulls == numCols || activeRow == numRows-1 {
		showSolution()
		activeRow = -1
	} else {
		activeRow++
		guess = 0
	}
}

func showSolution() {
	pics := document.GetElementsByClassName(solutionClass)
	for _, p := range pics {
		p.SetAttribute(atom.Style.String(), "")
	}
}
