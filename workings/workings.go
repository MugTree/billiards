package main

import (
	"github.com/goforj/godump"
)

type pocket struct {
	position int
	value    int
	isHit    bool
}

func main() {

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

	godump.DumpJSON(pockets)

}
