package app

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/MugTree/billiards/lib"
)

type pocket struct {
	position int
	value    int
	isHit    bool
}

type shot struct {
	pocket int
	timing time.Duration
}

const TotalPockets int = 9

const (
	redPlayer  string = "red"
	bluePlayer string = "blue"
)
const (
	gamePhaseShoot string = "shoot"
	gamePhaseScore string = "score"
)

type gameState struct {
	CurrentPlayer   string         `json:"current_player"`
	RedTotal        int64          `json:"red_total"`
	BlueTotal       int64          `json:"blue_total"`
	Phase           string         `json:"phase"`
	PocketsHit      []int          `json:"pocket_numbers_hit"`
	UserCalculation int            `json:"scorer_sum"`
	ScorerError     bool           `json:"scorer_error"`
	Score           int            `json:"shot_sum"`
	Table           map[int]pocket `json:"pockets"`
}

func newGameState(player string) gameState {

	return gameState{
		PocketsHit:    []int{},
		CurrentPlayer: player,
		Phase:         gamePhaseShoot,
		Table:         initialisePockets(),
	}

}

func (gs gameState) AsSignals() (map[string]any, error) {
	return lib.StructToMap(gs)
}

func swapPlayer(state gameState) gameState {
	newPlayer := lib.If(state.CurrentPlayer == bluePlayer, redPlayer, bluePlayer)

	state.CurrentPlayer = newPlayer
	state.Phase = gamePhaseShoot
	state.UserCalculation = 0
	state.Score = 0
	state.PocketsHit = []int{}
	return state
}

func playShot(state gameState) gameState {

	pocketsHit := shoot()

	// shot has been played the new phase is score regardless of outcome
	state.Phase = gamePhaseScore
	if len(pocketsHit) == 0 {
		state.Score = 0
		return state
	}

	score := 0
	for _, pocket := range pocketsHit {
		m, hasPocket := state.Table[pocket.pocket]
		if hasPocket {
			score += m.value
		}
	}

	state.PocketsHit = lib.Map(pocketsHit, func(s shot) int { return s.pocket })
	state.Score = score

	if len(state.PocketsHit) == 0 {
		return state
	}

	for _, v := range state.PocketsHit {
		val := state.Table[v]
		val.isHit = true
		state.Table[v] = val
	}

	return state
}

func updateScore(state gameState) (gameState, error) {

	switch state.CurrentPlayer {
	case bluePlayer:
		state.BlueTotal = state.BlueTotal + int64(state.Score)
	case redPlayer:
		state.RedTotal = state.RedTotal + int64(state.Score)
	default:
		return state, fmt.Errorf("player input is wrong: %s", state.CurrentPlayer)
	}

	state.Phase = gamePhaseShoot
	state.UserCalculation = 0
	state.PocketsHit = []int{}
	state.Score = 0

	return state, nil
}

func shoot() []shot {

	pocketsHit := make([]shot, TotalPockets)
	for i := 1; i <= TotalPockets; i++ {

		r := rand.IntN(4)

		// 25% chance of pot going in
		// --------------------------------------------------
		if r == 1 {

			// a few secongs
			// ----------------------------------------
			randSecs := rand.IntN(3)
			d := time.Second * time.Duration(randSecs)

			s := shot{
				pocket: i,
				timing: d, // we'll a
			}
			pocketsHit[i] = s
		}
	}
	return pocketsHit
}

func initialisePockets() map[int]pocket {
	pockets := map[int]pocket{}
	pockets[1] = pocket{position: 1, value: 10}
	pockets[2] = pocket{position: 2, value: 10}
	pockets[3] = pocket{position: 3, value: 20}
	pockets[4] = pocket{position: 4, value: 10}
	pockets[5] = pocket{position: 5, value: 10}
	pockets[6] = pocket{position: 6, value: 100}
	pockets[7] = pocket{position: 7, value: 50}
	pockets[8] = pocket{position: 8, value: 50}
	pockets[9] = pocket{position: 9, value: 200}
	return pockets
}
