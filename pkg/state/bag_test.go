package state

import (
	"math/rand"
	"testing"

	pb "github.com/abergmeier/blau/api/pb"
)

func TestAddToBag(t *testing.T) {
	b := &pb.Bag{
		Stones: make([]uint32, pb.Color_WHITE+1, pb.Color_WHITE+1),
	}

	AddToBag(b, pb.Color_YELLOW, 6)

	b.Stones[pb.Color_WHITE] = 0
	b.Stones[pb.Color_BLUE] = 0
	b.Stones[pb.Color_YELLOW] = 6
	b.Stones[pb.Color_RED] = 0

	AddToBag(b, pb.Color_YELLOW, 4)

	b.Stones[pb.Color_WHITE] = 0
	b.Stones[pb.Color_BLUE] = 0
	b.Stones[pb.Color_YELLOW] = 10
	b.Stones[pb.Color_RED] = 0
}

func TestTakeRandomFromBag(t *testing.T) {

	b := &pb.Bag{
		Stones: []uint32{
			1,  // Color_BLUE
			3,  // Color_YELLOW
			5,  // Color_RED
			7,  // Color_BLACK
			11, // Color_WHITE
		},
	}

	r := rand.New(rand.NewSource(7))

	takeCount := uint32(10)
	cs := takeRandomColorsFromBag(b, r, takeCount)

	for _, count := range cs {
		takeCount -= count
	}

	if takeCount > 0 {
		t.Errorf("Missing %v stones.", takeCount)
	}

	expected := ColorMap{0, 1, 1, 5, 3}

	if cs != expected {
		t.Errorf("Took %v. Expected %v", cs, expected)
	}
}

func TestCountStones(t *testing.T) {
	b := &pb.Bag{
		Stones: []uint32{
			1, // Color_BLUE
			2, // Color_YELLOW
			3, // Color_RED
			4, // Color_BLACK
			5, // Color_WHITE
		},
	}

	c := countStones(b)

	if c != 15 {
		t.Errorf("Invalid count %v. Expected 15", c)
	}
}

func TestFindRandomColorFromBag(t *testing.T) {

	b := &pb.Bag{
		Stones: []uint32{
			1,  // Color_BLUE
			3,  // Color_YELLOW
			5,  // Color_RED
			7,  // Color_BLACK
			11, // Color_WHITE
		},
	}

	r := rand.New(rand.NewSource(1))

	c := findRandomColorFromBag(b, r, 27)

	if c != pb.Color_BLACK {
		t.Errorf("Found wrong color %v. Expected %v.", c, pb.Color_BLACK)
	}
}

func TestSelectColorFromBag(t *testing.T) {

	b := &pb.Bag{
		Stones: []uint32{
			2, // Color_BLUE
			2, // Color_YELLOW
			2, // Color_RED
			2, // Color_BLACK
			2, // Color_WHITE
		},
	}

	testSelectColor(t, b, 0, pb.Color_BLUE)
	testSelectColor(t, b, 1, pb.Color_BLUE)
	testSelectColor(t, b, 2, pb.Color_YELLOW)
	testSelectColor(t, b, 3, pb.Color_YELLOW)
	testSelectColor(t, b, 4, pb.Color_RED)
	testSelectColor(t, b, 5, pb.Color_RED)
	testSelectColor(t, b, 6, pb.Color_BLACK)
	testSelectColor(t, b, 7, pb.Color_BLACK)
	testSelectColor(t, b, 8, pb.Color_WHITE)
	testSelectColor(t, b, 9, pb.Color_WHITE)
}

func testSelectColor(t *testing.T, b *pb.Bag, i uint32, expected pb.Color) {

	c := selectColorFromBag(b, i)

	if c == expected {
		return
	}

	t.Errorf("Selected wrong color %v. Expected %v", c, expected)
}
