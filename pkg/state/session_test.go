package state

import "testing"

func TestSessionsAdd(t *testing.T) {
	s := NewSession("1234")

	s.Add("1234")

	for _, player := range s.List() {
		if player != "1234" {
			t.Fatalf("Unexpected Player %v", player)
		}
	}

	s.Add("2345")

	for _, player := range s.List() {
		if player != "1234" && player != "2345" {
			t.Fatalf("Unexpected Player %v", player)
		}
	}
}

func TestSessionsRemove(t *testing.T) {
	s := NewSession("1234")

	s.Remove("1234")

	for player := range s.List() {
		t.Errorf("Unexpected List result %v", player)
	}
}

func TestSessionsRemoveOwner(t *testing.T) {
	s := NewSession("1234")

	s.Add("2345")
	s.Remove("1234")

	o := s.Owner()
	if o != "2345" {
		t.Errorf("Unexpected owner %v", o)
	}
}

func TestSessionClose(t *testing.T) {
	s := NewSession("1234")

	s.Add("2345")

	s.Close()
	if s.Owner() != "" {
		t.Error("Owner set after Close")
	}
}
