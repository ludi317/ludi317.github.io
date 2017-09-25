package main

import (
	"strconv"

	"github.com/shurcooL/htmlg"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func render() *html.Node {

	table := &html.Node{Data: atom.Table.String(), Type: html.ElementNode}

	// Title
	table.AppendChild(htmlg.Text("Mastermind"))

	// Solution
	tr := htmlg.TR()
	for _, s := range solution {
		img := &html.Node{
			Type: html.ElementNode, Data: atom.Img.String(),
			Attr: []html.Attribute{
				{Key: atom.Src.String(), Val: imageDir + "color_" + strconv.Itoa(s) + ".gif"},
				{Key: atom.Height.String(), Val: "36"},
				{Key: atom.Width.String(), Val: "36"},
				{Key: atom.Class.String(), Val: solutionClass},
				{Key: atom.Style.String(), Val: "visibility: hidden;"},
			},
		}
		tr.AppendChild(htmlg.TD(img))
	}
	table.AppendChild(tr)

	// Trial rows
	trials := &html.Node{
		Data: atom.Table.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Id.String(), Val: gameTableID},
			{Key: "bgcolor", Val: "#9e9e9e"},
			{Key: atom.Style.String(), Val: "border: 3px black solid;"},
		},
	}
	for row := NUM_ROWS - 1; row >= 0; row-- {
		tr := htmlg.TR()
		sRow := strconv.Itoa(row)
		for col := 0; col < NUM_COLS; col++ {
			sCol := strconv.Itoa(col)
			img := &html.Node{
				Type: html.ElementNode, Data: atom.Img.String(),
				Attr: []html.Attribute{
					{Key: atom.Src.String(), Val: imageDir + "hole.gif"},
					{Key: atom.Height.String(), Val: "36"},
					{Key: atom.Width.String(), Val: "36"},
					{Key: atom.Onclick.String(), Val: "placeColor(" + sRow + "," + sCol + ")"},
					{Key: atom.Id.String(), Val: sRow + "-" + sCol},
				},
			}
			tr.AppendChild(htmlg.TD(img))
		}
		tr.AppendChild(renderGrader(row))
		trials.AppendChild(tr)
	}
	table.AppendChild(trials)

	// Color picker
	var tds []*html.Node
	for i := 1; i <= NUM_COLORS; i++ {
		s := strconv.Itoa(i)
		img := &html.Node{
			Type: html.ElementNode,
			Data: atom.Img.String(),
			Attr: []html.Attribute{
				{Key: atom.Src.String(), Val: imageDir + "color_" + s + ".gif"},
				{Key: atom.Onclick.String(), Val: "pickColor(" + s + ")"},
			},
		}
		tds = append(tds, htmlg.TD(img))
	}
	tr = htmlg.TR(htmlg.TD(TBL(htmlg.TR(tds...))))
	table.AppendChild(tr)

	// Buttons
	solver := &html.Node{
		Data: atom.Input.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Onclick.String(), Val: "solve()"},
			{Key: atom.Value.String(), Val: "I'm Feeling Lazy"},
			{Key: atom.Type.String(), Val: "button"},
		},
	}

	newGame := &html.Node{
		Data: atom.Input.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Onclick.String(), Val: "reload()"},
			{Key: atom.Value.String(), Val: "New Game"},
			{Key: atom.Type.String(), Val: "button"},
		},
	}

	table.AppendChild(htmlg.TR(TBL(
		htmlg.TR(htmlg.TD(htmlg.TD(solver))),
		htmlg.TR(htmlg.TD(htmlg.TD(newGame))),
	)))

	return &html.Node{
		Data: atom.Div.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: atom.Style.String(), Val: `text-align: center; margin-top: 50px;`}},
		FirstChild: &html.Node{
			Data:       atom.Span.String(),
			Type:       html.ElementNode,
			Attr:       []html.Attribute{{Key: atom.Style.String(), Val: `display: inline-block; margin-left: 30px; margin-right: 30px;`}},
			FirstChild: table,
		},
	}
}

func renderGrader(row int) *html.Node {
	trs := []*html.Node{}
	for i := 0; i < 2; i++ {
		tr := htmlg.TR()
		for j := 0; j < NUM_COLS/2; j++ {
			img := &html.Node{
				Type: html.ElementNode, Data: atom.Img.String(),
				Attr: []html.Attribute{
					{Key: atom.Src.String(), Val: ""},
					{Key: atom.Height.String(), Val: "14"},
					{Key: atom.Width.String(), Val: "14"},
					{Key: atom.Class.String(), Val: "graderRow" + strconv.Itoa(row)},
				},
			}
			tr.AppendChild(htmlg.TD(img))
		}
		trs = append(trs, tr)
	}
	return htmlg.TD(TBL(trs...))
}

func TBL(nodes ...*html.Node) *html.Node {
	t := &html.Node{
		Type: html.ElementNode, Data: atom.Table.String(),
	}
	for _, n := range nodes {
		t.AppendChild(n)
	}
	return t
}
