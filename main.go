package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	WIDTH  = 4
	HEIGHT = 4
)

type Universe struct {
	width  int
	height int
	grid   [WIDTH][HEIGHT]bool
}

func (uni *Universe) New(width, height int) {
	(*uni).width = width
	(*uni).height = height
}

func (uni Universe) Show() {
	fmt.Println("Universe:")
	for y := 0; y < uni.height; y++ {
		for x := 0; x < uni.width; x++ {
			if uni.grid[y][x] {
				fmt.Print("*")
			} else {
				fmt.Print("#")
			}
		}

		fmt.Println()
	}
}

func (uni *Universe) Set(x, y int, isAlive bool) {
	(*uni).grid[y][x] = isAlive
}

func (uni *Universe) Init(width, height int) {
	(*uni).New(width, height)

	for i := 0; i < (uni.width * uni.height / 4); i++ {
		uni.Set(rand.Intn(uni.width), rand.Intn(uni.height), true)
	}
}

func (uni Universe) IsAlive(x, y int) bool {
	x = (x + uni.width) % uni.width
	y = (y + uni.height) % uni.height
	return uni.grid[y][x]
}

func (uni Universe) CountNeighbors(x, y int) (count int) {
	for row := y - 1; row <= y+1; row++ {
		for column := x - 1; column <= x+1; column++ {
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

func (uni1 *Universe) CopyData(uni2 Universe) {
	(*uni1).height = uni2.height
	(*uni1).width = uni2.width
	for y := 0; y < uni1.height; y++ {
		for x := 0; x < uni1.width; x++ {
			uni1.Set(x, y, uni2.IsAlive(x, y))
		}
	}
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

func Game() {
	var uni Universe
	uni.Init(WIDTH, HEIGHT)

	for i := 0; i < 15; i++ {
		uni.Show()
		uni.NextStep()
		time.Sleep(time.Second / 30)
	}
}

func main() {
	Game()
}
