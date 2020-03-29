package state

import (
	"fmt"

	pb "github.com/abergmeier/blau/api/pb"
)

type S interface {
	Add(PlayerUUID) error

	List() []PlayerUUID

	Remove(PlayerUUID)

	Owner() PlayerUUID
	Close()
}

type session pb.Session

func NewSession(owner PlayerUUID) S {
	s := &session{
		OwnerUuid: owner,
	}
	s.Add(owner)
	return s
}

// Add adds a new Player to the collection
func (s *session) Add(player PlayerUUID) error {

	if len(s.PlayerUuids) == 4 {
		return fmt.Errorf("Cannot add more than 4 players")
	}
	for _, p := range s.PlayerUuids {
		if p == player {
			return nil
		}
	}

	s.PlayerUuids = append(s.PlayerUuids, string(player))
	return nil
}

func (s *session) List() []PlayerUUID {
	return s.PlayerUuids
}

// Remove removes the Player from the collection
// If the owner is removed another Player becomes owner
func (s *session) Remove(player PlayerUUID) {

	for i, p := range s.PlayerUuids {
		if p == player {
			s.PlayerUuids = append(s.PlayerUuids[:i], s.PlayerUuids[i+1:]...)
			break
		}
	}

	if player != s.OwnerUuid {
		return
	}

	// We ought to find a new Owner
	if len(s.PlayerUuids) == 0 {
		s.Close()
		return
	}

	s.OwnerUuid = s.PlayerUuids[0]
}

func (s *session) Owner() PlayerUUID {
	return s.OwnerUuid
}

// Close has to be implemented
func (s *session) Close() {
	s.OwnerUuid = ""
	// Clear
	s.PlayerUuids = s.PlayerUuids[:0]
}
