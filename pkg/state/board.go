package state

import (
	pb "github.com/abergmeier/blau/api/pb"
)

func addToParking(b *pb.Board, rowI int, c pb.Color, amount uint32) {
	r := b.Rows[rowI]
	r.Parking.Color = c
	r.Parking.Amount = amount
}

func removeAllFromParking(b *pb.Board, rowI int) uint32 {
	r := b.Rows[rowI]

	a := r.Parking.Amount
	r.Parking.Color = pb.Color_BLUE
	r.Parking.Amount = 0
	return a
}
