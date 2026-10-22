package app

import (
	"math/rand/v2"
	"sort"
	"time"

	"github.com/MugTree/billiards/lib"
)

type pocket struct {
	Position int  `json:"position"`
	Value    int  `json:"value"`
	IsHit    bool `json:"is_hit"`
}

type potted struct {
	Pocket int           `json:"pocket"`
	Timing time.Duration `json:"timing"`
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
	PocketsHit      []int          `json:"pocket_hit"`
	UserCalculation int            `json:"user_calculation"`
	ScorerError     bool           `json:"scorer_error"`
	Score           int            `json:"score"`
	Pockets         map[int]pocket `json:"table"`
}

func (gs gameState) TableMapToSl() []pocket {

	if gs.Pockets == nil {
		LogInfo("table map is nil??")
		return []pocket{}
	}

	sl := make([]pocket, len(gs.Pockets))
	var i = 0
	for _, v := range gs.Pockets {
		sl[i] = v
		i++
	}

	sort.Slice(sl, func(i, j int) bool {
		return sl[i].Position < sl[j].Position
	})

	return sl
}

func newGameState(player string) gameState {

	return gameState{
		PocketsHit:    []int{},
		CurrentPlayer: player,
		Phase:         gamePhaseShoot,
		Pockets:       setupTable(),
	}
}

func (gs gameState) AsSignals() (map[string]any, error) {
	return lib.StructToMap(gs)
}

func playShot() []potted {

	pots := []potted{}

	for i := 1; i <= TotalPockets; i++ {

		r := rand.IntN(4)

		// 25% chance of pot going in
		// --------------------------------------------------
		if r == 1 {

			// a few secongs
			// ----------------------------------------
			randSecs := rand.IntN(3)
			d := time.Second * time.Duration(randSecs)

			s := potted{
				Pocket: i,
				Timing: d, // we'll a
			}
			pots = append(pots, s)
		}
	}
	return pots
}

func setupTable() map[int]pocket {
	pockets := map[int]pocket{}
	pockets[1] = pocket{Position: 1, Value: 10}
	pockets[2] = pocket{Position: 2, Value: 10}
	pockets[3] = pocket{Position: 3, Value: 20}
	pockets[4] = pocket{Position: 4, Value: 10}
	pockets[5] = pocket{Position: 5, Value: 10}
	pockets[6] = pocket{Position: 6, Value: 100}
	pockets[7] = pocket{Position: 7, Value: 50}
	pockets[8] = pocket{Position: 8, Value: 50}
	pockets[9] = pocket{Position: 9, Value: 200}
	return pockets
}
