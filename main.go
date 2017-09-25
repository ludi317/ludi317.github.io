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
	guess         = make([]int, NUM_COLS)
	solution      = make([]int, NUM_COLS)

	allCandidates []candidate
	feedbacks     []feedback
)

type candidate struct {
	code []int
	score int
}

type feedback struct {
	guess []int
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

func run() {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range solution {
		solution[i] = rand1.Intn(NUM_COLORS) + 1
	}
	document.Body().SetInnerHTML(htmlg.Render(render()))
}

func pickColor(s int) {
	selectedColor = s
	document.GetElementByID(gameTableID).SetAttribute(atom.Style.String(), `cursor: url("`+imageDir+`color_`+
		strconv.Itoa(selectedColor)+`.gif"), auto; border: 3px black solid;`)
}

func placeColor(row, col int) {
	if selectedColor != -1 && row == activeRow {
		document.GetElementByID(strconv.Itoa(row)+"-"+strconv.Itoa(col)).SetAttribute(atom.Src.String(),
			imageDir+"color_"+strconv.Itoa(selectedColor)+".gif")
		guess[col] = selectedColor
		grade()
	}
}

func solve() {
	go func() {
		allCandidates = genAllCandidates(NUM_COLS)
		for activeRow != -1 {
			for i, g := range knuthGuess() {
				selectedColor = g
				placeColor(activeRow, i)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func reload() {
	document.Location().Call("reload")
}

func grade() {
	for _, g := range guess {
		if g == 0 {
			return
		}
	}
	bulls, cows := score(guess, solution)
	feedbacks = append(feedbacks, feedback{guess, bulls, cows})
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
		guess = make([]int, NUM_COLS)
	}
}

func showSolution() {
	pics := document.GetElementsByClassName(solutionClass)
	for _, p := range pics {
		p.SetAttribute(atom.Style.String(), "")
	}
}
