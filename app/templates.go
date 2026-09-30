package app

import (
	"fmt"
	"strconv"

	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
}

func TemplateLayout(props PageProps, children ...Node) Node {

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Script(Src("/public/js/datastar.js"), Type("module")),
			// Link(Rel("stylesheet"), Href("/public/css/app.css")),
			Link(Rel("stylesheet"), Href("/public/css/main.css")),
		},
		Body: Group(children),
	})

}

func buildGameBoard(gs GameState) (Node, error) {

	sigs, err := gs.AsSignals()
	if err != nil {
		return nil, err
	}

	return TemplateGameBoard(gs, sigs), nil

}

func TemplateGameBoard(s GameState, sigs map[string]any) Node {

	pockets := pocketsMapAsSlice(s.GetTableState())

	return Div(ID("table"), ds.Signals(sigs),
		Div(
			ID("controls"),
			If(s.Phase == gamePhaseShoot,
				A(
					ds.On("click", "@get('/shoot')"),
					Style("text-decoration: underline"),
					Text("Start game"),
				)),
		),
		DrawBaize(pockets),
		If(s.Phase == gamePhaseScore, Div(
			Input(
				Name("scorer_sum"),
				ds.Bind("scorer_sum"),
			),
			Button(
				Text("Score"),
				ds.On("click", "@post('/score')"),
			),
		)),
	)
}

func DrawBaize(pockets []Pocket) Node {
	return Div(
		ID("baize"),
		Class("baize"),
		Map(pockets, func(p Pocket) Node {
			return Div(Classes{fmt.Sprintf("p%v", p.Position): true, "sunk": p.IsHit},
				Text(strconv.Itoa(p.Value)),
			)
		}),
	)
}
