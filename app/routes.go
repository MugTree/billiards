package app

import (
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/MugTree/billiards/app/db"

	"github.com/go-chi/chi/v5"
	"github.com/goforj/godump"
	"github.com/starfederation/datastar-go/datastar"
	"maragu.dev/gomponents"
)

func getRouter(_ *db.Queries) chi.Router {
	r := chi.NewRouter()
	//r.Use(httpDebugRequest)
	r.Handle("/public/*", httpNeuterDirectory(http.FileServer(http.FS(staticFS))))

	/**
	Start the game
	----------------------------------------------
	*/
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		players := []string{redPlayer, bluePlayer}
		startingPlayer := players[rand.IntN(2)]

		state := newGameState(startingPlayer)

		board, err := buildGameBoard(state)
		if err != nil {
			_httpLogAndError(w, r, err.Error())
			return
		}

		appLayout(
			pageProps{Title: "Home"},
			board,
		).Render(w)

	})

	/**
	Play a shot
	------------------------------------------------
	*/
	r.Get("/shoot", func(w http.ResponseWriter, r *http.Request) {

		state, ok := _gameStateFromSignalsOrError(w, r)
		if !ok {
			return
		}

		boardAfterShot, ok := _buildBoardOrError(w, r, playShot(state))
		if !ok {
			return
		}

		sse := datastar.NewSSE(w, r)
		sse.PatchElementGostar(boardAfterShot)
	})

	/**
	Score the shot
	-----------------------------------------------------
	*/
	r.Post("/score", func(w http.ResponseWriter, r *http.Request) {

		state, ok := _gameStateFromSignalsOrError(w, r)
		if !ok {
			return
		}

		fmt.Println("phase:", state.Phase)
		sse := datastar.NewSSE(w, r)

		// scoring error
		if state.UserCalculation != state.Score {
			godump.Dump("scoring error", state)

			state.ScorerError = true

			errorBoard, ok := _buildBoardOrError(w, r, state)
			if !ok {
				return
			}

			sse.PatchElementGostar(errorBoard)
			return
		}

		// swap player no score
		if state.UserCalculation == 0 {
			godump.Dump("no score", state)

			noScoreBoard, ok := _buildBoardOrError(w, r, swapPlayer(state))
			if !ok {
				return
			}

			sse.PatchElementGostar(noScoreBoard)
			return
		}

		godump.Dump("ALL GOOD", state)

		// all good
		state, err := updateScore(state)
		if err != nil {
			_httpLogAndError(w, r, err.Error())
			return
		}

		continueBoard, ok := _buildBoardOrError(w, r, state)
		if !ok {
			return
		}

		godump.Dump("NEW BOARD", state)
		sse.PatchElementGostar(continueBoard)
	})

	return r
}

// Helpers
// -----------------------------------------------------------------

func _gameStateFromSignalsOrError(w http.ResponseWriter, r *http.Request) (gameState, bool) {
	gs := gameState{}
	if err := datastar.ReadSignals(r, &gs); err != nil {
		_httpLogAndError(w, r, err.Error())
		return gs, false
	}

	return gs, true
}

func _buildBoardOrError(w http.ResponseWriter, r *http.Request, gs gameState) (gomponents.Node, bool) {
	board, err := buildGameBoard(gs)
	if err != nil {
		_httpLogAndError(w, r, err.Error())
		return nil, false
	}
	return board, true
}
