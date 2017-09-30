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
	d := &html.Node{
		Data: atom.Div.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Style.String(), Val: "font-size: 30px; margin: 10% auto; color: darkred;"},
		},
		FirstChild: htmlg.Text("Play Mastermind"),
	}

	table.AppendChild(d)

	// Solution
	tr := htmlg.TR()
	pieces := make([]int, numCols)

	for i, copySol := numCols-1, solution; i >= 0; i, copySol = i-1, copySol/10 {
		pieces[i] = copySol % 10
	}
	for _, s := range pieces {
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
	for row := numRows - 1; row >= 0; row-- {
		tr := htmlg.TR()
		sRow := strconv.Itoa(row)
		for col := 0; col < numCols; col++ {
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
	for i := 1; i <= numColors; i++ {
		s := strconv.Itoa(i)
		img := &html.Node{
			Type: html.ElementNode,
			Data: atom.Img.String(),
			Attr: []html.Attribute{
				{Key: atom.Src.String(), Val: imageDir + "color_" + s + ".gif"},
				{Key: atom.Onclick.String(), Val: "pickColor(" + s + ")"},
				{Key: atom.Onmouseover.String(), Val: "this.style.backgroundColor = '#fedac3'"},
				{Key: atom.Onmouseout.String(), Val: "this.style.backgroundColor = '#fefefe'"},
			},
		}
		tds = append(tds, htmlg.TD(img))
	}
	table4 := &html.Node{
		Type: html.ElementNode, Data: atom.Table.String(),
		Attr: []html.Attribute{
			{Key: atom.Id.String(), Val: "colorPickerID"},
		},
	}
	table4.AppendChild(htmlg.TR(tds...))
	tr = htmlg.TR(htmlg.TD(table4))
	table.AppendChild(tr)

	// Buttons
	solver := &html.Node{
		Data: atom.Input.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Onclick.String(), Val: "solve()"},
			{Key: atom.Value.String(), Val: "Solve in 5 or less moves"},
			{Key: atom.Type.String(), Val: "button"},
			{Key: atom.Style.String(), Val: `font-size: large;`},
		},
	}

	newGame := &html.Node{
		Data: atom.Input.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{
			{Key: atom.Onclick.String(), Val: "reload()"},
			{Key: atom.Value.String(), Val: "New Game"},
			{Key: atom.Type.String(), Val: "button"},
			{Key: atom.Style.String(), Val: `font-size: large;`},
		},
	}

	table.AppendChild(htmlg.TR(TBL(
		htmlg.TR(htmlg.TD(htmlg.TD(solver))),
		htmlg.TR(htmlg.TD(htmlg.TD(newGame))),
	)))

	return &html.Node{
		Data: atom.Div.String(),
		Type: html.ElementNode,
		Attr: []html.Attribute{{Key: atom.Style.String(), Val: `text-align: center;`}},
		FirstChild: &html.Node{
			Data:       atom.Div.String(),
			Type:       html.ElementNode,
			Attr:       []html.Attribute{{Key: atom.Style.String(), Val: `display: inline-block;`}},
			FirstChild: table,
		},
	}
}

func renderGrader(row int) *html.Node {
	trs := []*html.Node{}
	for i := 0; i < 2; i++ {
		tr := htmlg.TR()
		for j := 0; j < numCols/2; j++ {
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
