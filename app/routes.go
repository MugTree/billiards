package app

import (
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/MugTree/billiards/app/db"
	"github.com/MugTree/billiards/lib"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

func isValidInput(w http.ResponseWriter, r *http.Request) (gameState, bool) {

	var game = gameState{}

	if err := datastar.ReadSignals(r, &game); err != nil {
		_httpLogAndError(w, r, err.Error())
		return game, false
	}

	// Add more validation when required
	// -----------------------------------------------

	// is the player value correct?
	//
	if game.CurrentPlayer != redPlayer && game.CurrentPlayer != bluePlayer {
		_httpLogAndError(w, r, fmt.Errorf("Validation error: player needs to be red or blue").Error())
		return game, false
	}

	return game, true

}

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
		board := gameBoard(newGameState(startingPlayer))

		appLayout(
			pageProps{Title: "Home"},
			board,
		).Render(w)

	})

	//
	// Play a shot
	//------------------------------------------------
	r.Get("/shoot", func(w http.ResponseWriter, r *http.Request) {

		game, ok := isValidInput(w, r)
		if !ok {
			return
		}

		pocketsHit := playShot()

		score := 0
		for _, p := range pocketsHit {
			m, _ := game.Pockets[p.Pocket]
			score += m.Value
		}
		game.Score = score

		game.Phase = gamePhaseScore
		game.PocketsHit = lib.Map(pocketsHit, func(s potted) int { return s.Pocket })

		for _, p := range game.PocketsHit {
			pocket := game.Pockets[p]
			pocket.IsHit = true
			game.Pockets[p] = pocket
		}

		sse := datastar.NewSSE(w, r)
		sse.PatchElementGostar(gameBoard(game))

	})

	// Score the shot
	// -----------------------------------------------------
	r.Post("/score", func(w http.ResponseWriter, r *http.Request) {

		game, ok := isValidInput(w, r)
		if !ok {
			return
		}

		sse := datastar.NewSSE(w, r)

		// scorer error
		// ------------------------------------------------------
		if game.UserCalculation != game.Score {

			game.ScorerError = true
			sse.PatchElementGostar(gameBoard(game))
			return
		}

		// no score change player
		// ---------------------------------------------
		if game.Score == 0 {

			game.CurrentPlayer = lib.If(game.CurrentPlayer == bluePlayer, redPlayer, bluePlayer)
			game.Phase = gamePhaseShoot
			game.UserCalculation = 0
			game.Score = 0
			game.PocketsHit = []int{}
			sse.PatchElementGostar(gameBoard(game))
			return
		}

		// build a break
		// --------------------------------------------
		if game.CurrentPlayer == bluePlayer {
			game.BlueTotal = game.BlueTotal + int64(game.Score)
		} else {
			game.RedTotal = game.RedTotal + int64(game.Score)
		}

		game.Phase = gamePhaseShoot
		game.UserCalculation = 0
		game.PocketsHit = []int{}
		game.Score = 0
		sse.PatchElementGostar(gameBoard(game))
	})

	return r
}
