package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/shurcooL/go/gopherjs_http/jsutil"
	"github.com/shurcooL/htmlg"
	"honnef.co/go/js/dom"
)

const (
	numRows   = 10
	numCols   = 4
	numColors = 6

	gameTableID   = "gameTableID"
	solutionClass = "solutionClass"
	imageDir      = "images/"
)

var (
	document      dom.HTMLDocument
	selectedColor int
	activeRow     int
	solution      int

	feedbackHash int
	guess        int
)

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
	solution = generateSolution()
	println(solution)
	document.Body().SetInnerHTML(htmlg.Render(render()))
}
