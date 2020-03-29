package state

import (
	"math/rand"

	pb "github.com/abergmeier/blau/api/pb"
)

func AddToBag(bag *pb.Bag, c pb.Color, amount uint32) {
	bag.Stones[c] += amount
}

type ColorMap [pb.Color_WHITE + 1]uint32

func TakeRandomColorsFromBag(bag *pb.Bag, amount uint32) ColorMap {

	r := rand.New(rand.NewSource(3))
	return takeRandomColorsFromBag(bag, r, amount)
}

func takeRandomColorsFromBag(bag *pb.Bag, r *rand.Rand, amount uint32) ColorMap {

	toTake := ColorMap{}

	stoneCount := countStones(bag)

	for i := uint32(0); i != amount; i++ {
		c := findRandomColorFromBag(bag, r, stoneCount)
		toTake[c]++
	}

	for c, count := range toTake {
		bag.Stones[c] -= count
	}

	return toTake
}

func countStones(bag *pb.Bag) uint32 {
	count := uint32(0)
	for _, c := range bag.Stones {
		count += c
	}

	return count
}

func findRandomColorFromBag(bag *pb.Bag, r *rand.Rand, max uint32) pb.Color {

	stoneIndex := uint32(r.Intn(int(max)))
	return selectColorFromBag(bag, stoneIndex)
}

func selectColorFromBag(bag *pb.Bag, stoneIndex uint32) pb.Color {
	for c, count := range bag.Stones {
		if stoneIndex < count {
			return pb.Color(c)
		}

		stoneIndex -= count
	}

	panic("Unexpected")
}
