package app

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/MugTree/billiards/app/db"
	"github.com/MugTree/billiards/lib"
	"github.com/goforj/godump"

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
		startBoard := gameBoard(newGameState(startingPlayer))

		appLayout(
			pageProps{Title: "Home"},
			startBoard,
		).Render(w)

	})

	//
	// Play a shot
	//------------------------------------------------
	r.Get("/shoot", func(w http.ResponseWriter, r *http.Request) {

		initialState, ok := isValidInput(w, r)
		if !ok {
			return
		}

		// simple game atm assumes that no balls hang about
		// -------------------------------------------------

		initialState.Pockets = setupTable()

		sse := datastar.NewSSE(w, r)

		if initialState.IsBreak {

			// redraw board
			sse.PatchElementGostar(gameBoard(initialState))
		}

		afterShot := func(game gameState) gameState {

			for i := 1; i <= len(game.Pockets); i++ {

				r := rand.IntN(4)

				// 25% chance of pot going in??? :/
				// --------------------------------------------------
				if r == 1 {

					pocket := game.Pockets[i]
					pocket.TimeToDrop = rand.IntN(2) + 1
					pocket.IsHit = true

					game.Pockets[i] = pocket
					game.Score += pocket.Value
					game.IsBreak = true

				}
			}

			game.Phase = gamePhaseScore
			return game

		}(initialState)

		//sse := datastar.NewSSE(w, r)

		fmt.Printf("There are %v balls potted...", afterShot.GetPottedCount())
		fmt.Println()

		for _, v := range afterShot.Pockets {

			if v.IsHit {

				dropTime := v.TimeToDrop * int(time.Second)
				time.Sleep(time.Duration(dropTime))
				sse.PatchElementGostar(drawPocket(v))
				fmt.Println("animating...")
				fmt.Println()

			}
		}

		sse.PatchElementGostar(gameBoard(afterShot))

		b, err := afterShot.ToJSON()
		if err != nil {
			_httpLogAndError(w, r, err.Error())
			return
		}

		sse.PatchSignals(b)

		debugObj(afterShot.Phase)
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
			game.Pockets = setupTable()
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
		game.Pockets = setupTable()
		game.Score = 0
		sse.PatchElementGostar(gameBoard(game))
	})

	return r
}

func debugObj(obj any) {

	godump.DumpJSON(obj)
	//fmt.Printf("%+v\n", obj)

}
