package game

func (uni Universe) ShowHTML() (output string) {
	for y := 0; y < uni.height; y++ {
		for x := 0; x < uni.width; x++ {
			if uni.grid[y][x] {
				output += ALIVE
			} else {
				output += DEAD
			}
		}
		output += "\n"
	}

	return
}
