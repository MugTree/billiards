package main

import (
	"math/rand/v2"

	"github.com/goforj/godump"
)

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

func main() {

	s := shoot()

	godump.Dump(len(s), s)

}
