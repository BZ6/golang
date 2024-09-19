package game

import (
	"fmt"
)

func (uni Universe) ShowConsole() {
	fmt.Println("Universe:")
	
	for y := 0; y < uni.height; y++ {
		for x := 0; x < uni.width; x++ {
			if uni.grid[y][x] {
				fmt.Print(ALIVE)
			} else {
				fmt.Print(DEAD)
			}
		}

		fmt.Println()
	}
}
