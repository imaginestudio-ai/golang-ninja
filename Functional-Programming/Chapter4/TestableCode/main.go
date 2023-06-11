package main

import (
	"math/rand"

	"github.com/ImagineDevOps DevOps/Chapter4/TestableCode/player"
)

func main() {
	random := rand.Intn(2)
	player.PlayerSelectPure(random)
	// start the game
}
