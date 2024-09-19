package game

func (uni Universe) IsAlive(x, y int) bool {
	x = (x + uni.width) % uni.width
	y = (y + uni.height) % uni.height
	return uni.grid[y][x]
}

func (uni Universe) CountNeighbors(x, y int) (count int) {
	for row := y - 1; row <= y + 1; row++ {
		for column := x - 1; column <= x + 1; column++ {
			if uni.IsAlive(column, row) {
				count++
			}
		}
	}

	if uni.IsAlive(x, y) {
		count--
	}
	return
}

func (uni Universe) NextLive(x, y int) bool {
	count := uni.CountNeighbors(x, y)

	if !uni.IsAlive(x, y) && count == 3 {
		return true
	}
	if uni.IsAlive(x, y) && 2 <= count && count <= 3 {
		return true
	}

	return false
}

func (uni *Universe) NextStep() {
	var uniNext Universe
	uniNext.New(uni.width, uni.height)

	for y := 0; y < uni.height; y++ {
		for x := 0; x < uni.width; x++ {
			uniNext.Set(x, y, uni.NextLive(x, y))
		}
	}

	uni.CopyData(uniNext)
}
