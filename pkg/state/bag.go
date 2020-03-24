package state

import (
	"github.com/abergmeier/blau/pkg/api/pb"
)

func AddToBag(bag *pb.Bag, c pb.Color, amount uint16) {
	bag.Stones.Set(c, bag.Stones.Get(c)+amount)
}

func TakeFromBag(bag *pb.Bag, c pb.Color, amount uint16) {
	cStones, ok := bag.Stones.get(c)
	if !ok {
		panic("Stones are fucked up")
	}

	cStones -= amount
}
