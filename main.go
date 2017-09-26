package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/shurcooL/go/gopherjs_http/jsutil"
	"github.com/shurcooL/htmlg"
	"golang.org/x/net/html/atom"
	"honnef.co/go/js/dom"
)

const (
	NUM_ROWS   = 10
	NUM_COLS   = 4
	NUM_COLORS = 6

	gameTableID   = "gameTableID"
	solutionClass = "solutionClass"
	imageDir      = "images/"
)

var (
	document      dom.HTMLDocument
	selectedColor int
	activeRow     int
	solution      int

	//allCandidates  = genAllCandidates(NUM_COLS)
	allCandidates = []int{}
	feedbackHash  int
	guess         int
)

type feedback struct {
	guess int
	bulls int
	cows  int
}

func main() {
	js.Global.Set("pickColor", jsutil.Wrap(pickColor))
	js.Global.Set("placeColor", jsutil.Wrap(placeColor))
	js.Global.Set("solve", jsutil.Wrap(solve))
	js.Global.Set("reload", jsutil.Wrap(reload))
	document = dom.GetWindow().Document().(dom.HTMLDocument)
	document.AddEventListener("DOMContentLoaded", false, func(dom.Event) {
		go run()
	})
}

func generateSolution() int {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := 0
	for i := 0; i < NUM_COLS; i++ {
		s += rand1.Intn(NUM_COLORS) + 1
		if i != NUM_COLS-1 {
			s *= 10
		}
	}
	return s
}

func run() {
	solution = generateSolution()
	println(solution)
	document.Body().SetInnerHTML(htmlg.Render(render()))
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
	for i := 0; i < NUM_COLS-col-1; i++ {
		m *= 10
		copyGuess /= 10
	}
	curGuess := copyGuess % 10
	return guess + m*(selectedColor-curGuess)
}

func solve() {
	go func() {
		for activeRow != -1 {
			m := k.move
			for col := NUM_COLS - 1; col >= 0; col-- {
				selectedColor = m % 10
				m /= 10
				placeColor(activeRow, col)
			}
			k = k.bullsAndCows[feedbackHash]
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func reload() {
	document.Location().Call("reload")
}

func hasZeros(guess int) bool {
	for i := 0; i < NUM_COLS; i++ {
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

	bulls, cows := score(guess, solution)
	feedbackHash = hash(bulls, cows)
	pegHoles := document.GetElementsByClassName("graderRow" + strconv.Itoa(activeRow))
	i := 0
	for ; i < cows; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_6.gif")
	}
	for ; i < cows+bulls; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_5.gif")
	}

	if bulls == NUM_COLS || activeRow == NUM_ROWS-1 {
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
