package app

import (
	"fmt"
	"strconv"

	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type pageProps struct {
	Title       string
	Description string
}

func appLayout(props pageProps, children ...Node) Node {

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

func buildGameBoard(gs gameState) (Node, error) {

	sigs, err := gs.AsSignals()
	if err != nil {
		return nil, err
	}

	return gameBoard(gs, sigs), nil

}

func gameBoard(s gameState, sigs map[string]any) Node {

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
		drawBaize(pockets),
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

func drawBaize(pockets []pocket) Node {
	return Div(
		ID("baize"),
		Class("baize"),
		Map(pockets, func(p pocket) Node {
			return Div(Classes{fmt.Sprintf("p%v", p.position): true, "sunk": p.isHit},
				Text(strconv.Itoa(p.value)),
			)
		}),
	)
}
