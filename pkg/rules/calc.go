package rules

import (
	"github.com/abergmeier/blau/pkg/api/pb"
	"github.com/abergmeier/blau/pkg/state"
)

func CalcPoints(b *pb.Board, bag *pb.Bag) {
	for i, _ := b.Rows {
		calcRowPoints(b, i)
	}
}

func calcRowPoints(b *pb.Board, rowI int, bag *pb.Bag) {
	r := b.Rows[rowI]

	commitThreshold := rowI
	if r.Parking.Amount < commitThreshold {
		return
	}

	amount := state.RemoveAllFromParking(b, rowI)

	ci := colorIndex(r.Parking.Color, rowI int)
	setBit(&r.StoredBitfield, ci)
	amount -= 1
	state.AddToBag(bag, r.Parking.Color, amount)
}

func setBit(bitfield *uint32, bitIndex int) {
	bitfield = 1 >> bitIndex
}

func colorIndex(c pb.Color, rowI int) int {
	return c + rowI
}
