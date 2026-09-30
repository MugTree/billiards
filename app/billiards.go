package app

import (
	"fmt"
	"math/rand/v2"
	"sort"

	"github.com/MugTree/billiards/lib"
)

type pocket struct {
	position int
	value    int
	isHit    bool
}

const (
	redPlayer  string = "red"
	bluePlayer string = "blue"
)
const (
	gamePhaseShoot string = "shoot"
	gamePhaseScore string = "score"
)

type gameState struct {
	CurrentPlayer     string `json:"current_player"`
	RedTotal          int64  `json:"red_total"`
	BlueTotal         int64  `json:"blue_total"`
	Phase             string `json:"phase"`
	PocketsNumbersHit []int  `json:"pocket_numbers_hit"`
	ScorerSum         int    `json:"scorer_sum"`
	ScorerError       bool   `json:"scorer_error"`
	ShotSum           int    `json:"shot_sum"`
}

func startGame() gameState {

	players := []string{redPlayer, bluePlayer}
	startingPlayer := players[rand.IntN(2)]

	return gameState{
		PocketsNumbersHit: []int{},
		CurrentPlayer:     startingPlayer,
		Phase:             gamePhaseShoot,
	}
}

func (gs gameState) GetTableState() map[int]pocket {

	pm := getPocketMap()

	if len(gs.PocketsNumbersHit) == 0 {
		return pm
	}

	pocketsHit := gs.PocketsNumbersHit

	for _, v := range pocketsHit {
		val := pm[v]
		val.isHit = true
		pm[v] = val
	}

	return pm
}

func (gs gameState) AsSignals() (map[string]any, error) {
	return lib.StructToMap(gs)
}

func swapPlayer(gs gameState) gameState {
	newPlayer := lib.If(gs.CurrentPlayer == bluePlayer, redPlayer, bluePlayer)
	gs.CurrentPlayer = newPlayer
	gs.Phase = gamePhaseShoot
	gs.ScorerSum = 0
	gs.ShotSum = 0
	gs.PocketsNumbersHit = []int{}

	return gs
}

func playShot(gs gameState, shotResult []int) gameState {

	gs.Phase = gamePhaseScore
	gs.PocketsNumbersHit = shotResult

	if len(shotResult) == 0 {
		gs.ShotSum = 0
		return gs
	}

	pockets := getPocketMap()
	score := 0

	for _, pn := range gs.PocketsNumbersHit {
		m, ok := pockets[pn]
		// check against tampering
		if ok {
			score += m.value
		}
	}

	gs.ShotSum = score
	return gs
}

func updatePlayerScore(gs gameState) (gameState, error) {

	switch gs.CurrentPlayer {
	case bluePlayer:
		gs.BlueTotal = gs.BlueTotal + int64(gs.ShotSum)
	case redPlayer:
		gs.RedTotal = gs.RedTotal + int64(gs.ShotSum)
	default:
		return gs, fmt.Errorf("player input is wrong: %s", gs.CurrentPlayer)
	}

	gs.Phase = gamePhaseShoot
	gs.ScorerSum = 0
	gs.PocketsNumbersHit = []int{}
	gs.ShotSum = 0

	return gs, nil
}

func pocketsMapAsSlice(pm map[int]pocket) []pocket {
	pockets := []pocket{}
	for k, p := range pm {
		p.position = k
		pockets = append(pockets, p)
	}

	sort.Slice(pockets, func(i, j int) bool { return pockets[i].position < pockets[j].position })
	return pockets
}

// generate an slice of ints up to 9 long consisting of
// numbers 1..9
func shoot() (pocketsHit []int) {
	for i := 1; i < 10; i++ {
		r := rand.IntN(4)
		if r == 1 {
			pocketsHit = append(pocketsHit, i)
		}
	}
	return pocketsHit
}

func getPocketMap() map[int]pocket {
	pockets := map[int]pocket{}
	pockets[1] = pocket{value: 10}
	pockets[2] = pocket{value: 10}
	pockets[3] = pocket{value: 20}
	pockets[4] = pocket{value: 10}
	pockets[5] = pocket{value: 10}
	pockets[6] = pocket{value: 100}
	pockets[7] = pocket{value: 50}
	pockets[8] = pocket{value: 50}
	pockets[9] = pocket{value: 200}
	return pockets
}
