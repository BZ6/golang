package game

import (
	"time"
)

// Game for console
func GameConsole() {
	var uni Universe
	uni.Init(WIDTH, HEIGHT)

	for i := 0; i < REPEATS; i++ {
		uni.ShowConsole()
		uni.NextStep()
		time.Sleep(time.Second / FPS)
	}
}
