package app

import (
	"fmt"
	"math/rand/v2"
	"sort"

	"github.com/MugTree/billiards/lib"
)

type Pocket struct {
	Position int
	Value    int
	IsHit    bool
}

const (
	redPlayer  string = "red"
	bluePlayer string = "blue"
)
const (
	gamePhaseShoot string = "shoot"
	gamePhaseScore string = "score"
)

type GameState struct {
	CurrentPlayer     string `json:"current_player"`
	RedTotal          int64  `json:"red_total"`
	BlueTotal         int64  `json:"blue_total"`
	Phase             string `json:"phase"`
	PocketsNumbersHit []int  `json:"pocket_numbers_hit"`
	ScorerSum         int    `json:"scorer_sum"`
	ScorerError       bool   `json:"scorer_error"`
	ShotSum           int    `json:"shot_sum"`
}

func startGame() GameState {

	players := []string{redPlayer, bluePlayer}
	startingPlayer := players[rand.IntN(2)]

	return GameState{
		PocketsNumbersHit: []int{},
		CurrentPlayer:     startingPlayer,
		Phase:             gamePhaseShoot,
	}
}

func swapPlayer(gs GameState) GameState {
	newPlayer := lib.If(gs.CurrentPlayer == bluePlayer, redPlayer, bluePlayer)
	gs.CurrentPlayer = newPlayer
	gs.Phase = gamePhaseShoot
	gs.ScorerSum = 0
	gs.ShotSum = 0
	gs.PocketsNumbersHit = []int{}

	return gs
}

func playShot(gs GameState) GameState {

	gs.Phase = gamePhaseScore

	shot := shoot()
	gs.PocketsNumbersHit = shot

	if len(shot) == 0 {
		gs.ShotSum = 0
		return gs
	}

	pockets := getPocketMap()
	score := 0

	for _, pn := range gs.PocketsNumbersHit {
		m, ok := pockets[pn]
		// check against tampering
		if ok {
			score += m.Value
		}
	}

	gs.ShotSum = score
	return gs
}

func updatePlayerScore(gs GameState) (GameState, error) {

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

func (gs GameState) GetTableState() map[int]Pocket {

	pm := getPocketMap()

	if len(gs.PocketsNumbersHit) == 0 {
		return pm
	}

	pocketsHit := gs.PocketsNumbersHit

	for _, v := range pocketsHit {
		val := pm[v]
		val.IsHit = true
		pm[v] = val
	}

	return pm
}

func (gs GameState) AsSignals() (map[string]any, error) {
	return lib.StructToMap(gs)
}

func pocketsMapAsSlice(pm map[int]Pocket) []Pocket {
	pockets := []Pocket{}
	for k, p := range pm {
		p.Position = k
		pockets = append(pockets, p)
	}

	sort.Slice(pockets, func(i, j int) bool { return pockets[i].Position < pockets[j].Position })
	return pockets
}

// generate an slice of ints up to 9 long consisting of
// numbers 1..9
func shoot() []int {
	nums := []int{}
	for i := 1; i < 10; i++ {
		r := rand.IntN(4)
		if r == 1 {
			nums = append(nums, i)
		}
	}
	return nums
}

func getPocketMap() map[int]Pocket {
	pockets := map[int]Pocket{}
	pockets[1] = Pocket{Value: 10}
	pockets[2] = Pocket{Value: 10}
	pockets[3] = Pocket{Value: 20}
	pockets[4] = Pocket{Value: 10}
	pockets[5] = Pocket{Value: 10}
	pockets[6] = Pocket{Value: 100}
	pockets[7] = Pocket{Value: 50}
	pockets[8] = Pocket{Value: 50}
	pockets[9] = Pocket{Value: 200}
	return pockets
}
