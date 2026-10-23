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

func gameBoard(state gameState) Node {

	// has been validated at this point
	sigs, _ := state.AsSignals()

	return Div(ID("table"), ds.Signals(sigs),
		Div(
			ID("controls"),
			If(state.Phase == gamePhaseShoot,
				A(
					ds.On("click", "@get('/shoot')"),
					Style("text-decoration: underline"),
					Text("Play shot!"),
				)),
		),

		drawBaize(state),

		If(state.Phase == gamePhaseScore, Div(
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

func drawBaize(state gameState) Node {

	sl := state.TableMapToSl()

	return Div(ID("baize"), Class("baize"),

		Map(sl, func(p pocket) Node {
			return drawPocket(p)
		}),
	)
}

func drawPocket(p pocket) Node {
	return Div(
		ID(fmt.Sprintf("pocket-%v", p.Position)),
		Classes{fmt.Sprintf("p%v", p.Position): true, "sunk": p.IsHit},
		Text(strconv.Itoa(p.Value)),
	)
}
