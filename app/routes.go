package app

import (
	"fmt"
	"net/http"

	"github.com/MugTree/billiards/app/db"

	"github.com/go-chi/chi/v5"
	"github.com/goforj/godump"
	"github.com/starfederation/datastar-go/datastar"
	"maragu.dev/gomponents"
)

func getRouter(queries *db.Queries) chi.Router {
	r := chi.NewRouter()
	//r.Use(httpDebugRequest)
	r.Handle("/public/*", httpNeuterDirectory(http.FileServer(http.FS(staticFS))))
	r.Get("/", startHandler(queries))
	r.Get("/shoot", shootHandler)
	r.Post("/score", scoreHandler)
	return r
}

func startHandler(_ *db.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		state := startGame()
		board, err := buildGameBoard(state)
		if err != nil {
			httpLogAndError(w, r, err.Error())
			return
		}

		TemplateLayout(
			PageProps{Title: "Home"},
			board,
		).Render(w)

	}
}

func shootHandler(w http.ResponseWriter, r *http.Request) {

	gs, ok := _gameStateFromSignalsOrError(w, r)
	if !ok {
		return
	}

	boardAfterShot, ok := _buildBoardOrError(w, r, playShot(gs))
	if !ok {
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementGostar(boardAfterShot)

}

func scoreHandler(w http.ResponseWriter, r *http.Request) {

	gs, ok := _gameStateFromSignalsOrError(w, r)
	if !ok {
		return
	}

	fmt.Println("phase:", gs.Phase)
	sse := datastar.NewSSE(w, r)

	// scoring error
	if gs.ScorerSum != gs.ShotSum {
		godump.Dump("scoring error", gs)

		gs.ScorerError = true

		errorBoard, ok := _buildBoardOrError(w, r, gs)
		if !ok {
			return
		}

		sse.PatchElementGostar(errorBoard)
		return
	}

	// swap player no score
	if gs.ScorerSum == 0 {
		godump.Dump("no score", gs)

		noScoreBoard, ok := _buildBoardOrError(w, r, swapPlayer(gs))
		if !ok {
			return
		}

		sse.PatchElementGostar(noScoreBoard)
		return
	}

	godump.Dump("ALL GOOD", gs)

	// all good
	gs, err := updatePlayerScore(gs)
	if err != nil {
		httpLogAndError(w, r, err.Error())
		return
	}

	continueBoard, ok := _buildBoardOrError(w, r, gs)
	if !ok {
		return
	}

	godump.Dump("NEW BOARD", gs)
	sse.PatchElementGostar(continueBoard)
}

func _gameStateFromSignalsOrError(w http.ResponseWriter, r *http.Request) (GameState, bool) {
	gs := GameState{}
	if err := datastar.ReadSignals(r, &gs); err != nil {
		httpLogAndError(w, r, err.Error())
		return gs, false
	}

	return gs, true
}

func _buildBoardOrError(w http.ResponseWriter, r *http.Request, gs GameState) (gomponents.Node, bool) {
	board, err := buildGameBoard(gs)
	if err != nil {
		httpLogAndError(w, r, err.Error())
		return nil, false
	}
	return board, true
}
