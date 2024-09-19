package game

import (
	"math/rand"
)

func (uni *Universe) New(width, height int) {
	(*uni).width = width
	(*uni).height = height
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

func (uni1 *Universe) CopyData(uni2 Universe) {
	(*uni1).height = uni2.height
	(*uni1).width = uni2.width
	for y := 0; y < uni1.height; y++ {
		for x := 0; x < uni1.width; x++ {
			uni1.Set(x, y, uni2.IsAlive(x, y))
		}
	}
}
