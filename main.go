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
	NUM_COLORS = 8

	gameTableID   = "gameTableID"
	dataRowAttr   = "data-row"
	dataColAttr   = "data-col"
	solutionClass = "solutionClass"
	checkBoxID    = "checkboxID"
	imageDir      = "images/"
)

var (
	document      dom.HTMLDocument
	selectedColor int
	activeRow     int
	guess         [NUM_COLS]int
	solution      [NUM_COLS]int
)

func main() {
	js.Global.Set("pickColor", jsutil.Wrap(pickColor))
	js.Global.Set("placeColor", jsutil.Wrap(placeColor))
	document = dom.GetWindow().Document().(dom.HTMLDocument)
	document.AddEventListener("DOMContentLoaded", false, func(dom.Event) {
		go run()
	})
}

func run() {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	perms := rand1.Perm(NUM_COLORS)
	for i := range solution {
		solution[i] = perms[i] + 1
	}
	println(solution)
	document.Body().SetInnerHTML(htmlg.Render(render()))
}

func pickColor(s int) {
	selectedColor = s
	document.GetElementByID(gameTableID).SetAttribute(atom.Style.String(), `cursor: url("`+imageDir+`color_`+
		strconv.Itoa(selectedColor)+`.gif"), auto; border: 3px black solid;`)
}

func placeColor(this dom.HTMLElement) {
	row := this.GetAttribute(dataRowAttr)
	if selectedColor != -1 && row == strconv.Itoa(activeRow) {
		this.SetAttribute(atom.Src.String(), imageDir+"color_"+strconv.Itoa(selectedColor)+".gif")
		iCol, _ := strconv.Atoi(this.GetAttribute(dataColAttr))
		guess[iCol] = selectedColor
		score()
	}
}

func score() {
	for _, g := range guess {
		if g == 0 {
			return
		}
	}
	bulls, cows := 0, 0
	colors := [NUM_COLORS + 1]int{}
	for i, g := range guess {
		if solution[i] == g {
			bulls++
		} else {
			if colors[g] < 0 {
				cows++
			}
			if colors[solution[i]] > 0 {
				cows++
			}
			colors[g]++
			colors[solution[i]]--
		}
	}
	pegHoles := document.GetElementsByClassName("graderRow" + strconv.Itoa(activeRow))
	i := 0
	for ; i < cows; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_8.gif")
	}
	for ; i < cows+bulls; i++ {
		pegHoles[i].SetAttribute(atom.Src.String(), imageDir+"color_7.gif")
	}

	if bulls == NUM_COLS {
		showSolution()
		js.Global.Call("alert", "Congrats! You solved it.")
		activeRow = -1
	} else if activeRow == NUM_ROWS-1 {
		showSolution()
		js.Global.Call("alert", "Sorry, you're out of tries. See solution.")
		activeRow = -1
	} else {
		activeRow++
		guess = [NUM_COLS]int{}
	}
}

func showSolution() {
	pics := document.GetElementsByClassName(solutionClass)
	for _, p := range pics {
		p.SetAttribute(atom.Style.String(), "")
	}
}
