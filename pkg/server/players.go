package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	pb "github.com/abergmeier/blau/api/pb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	lastSeenDocPath = "Players/LastSeen"
	storeDocPath    = "Players/Info"
)

type PlayerInfo struct {
	player   *pb.Player
	lastSeen int64
}

type playerServer struct {
	client       *firestore.Client
	lastSeen     *firestore.DocumentRef
	store        *firestore.DocumentRef
	ticker       *time.Ticker
	deleted      map[string]bool
	deletedMutex sync.Mutex
	rg           sync.WaitGroup
	stopReload   chan bool
}

func NewPlayerServer(client *firestore.Client, updateDuration time.Duration) (*playerServer, error) {

	lastSeen, err := ensureDoc(client, lastSeenDocPath)
	if err != nil {
		return nil, err
	}

	store, err := ensureDoc(client, storeDocPath)
	if err != nil {
		return nil, err
	}

	p := &playerServer{
		client:     client,
		lastSeen:   lastSeen,
		store:      store,
		ticker:     time.NewTicker(updateDuration),
		deleted:    map[string]bool{},
		stopReload: make(chan bool),
	}

	go p.reload()
	return p, nil
}

func (p *playerServer) Close() {
	p.ticker.Stop()
	p.stopReload <- true
	p.rg.Wait()
}

func ensureDoc(client *firestore.Client, docPath string) (*firestore.DocumentRef, error) {
	doc := client.Doc(docPath)

	_, err := doc.Create(context.Background(), map[string]interface{}{})
	if err == nil {
		return doc, nil
	}

	// There may be an error we care about

	s, ok := status.FromError(err)
	if ok && s.Code() == codes.AlreadyExists {
		// Disregard error since successfully ensured
		return doc, nil
	}

	// There is an actual error
	return nil, err
}

func (p *playerServer) reload() {

	p.rg.Add(1)
	defer p.rg.Done()
	for {
		select {
		case <-p.stopReload:
			return
		case <-p.ticker.C:
			p.deleteTimedOut()
		}
	}
}

func (p *playerServer) deleteTimedOut() {
	doc, err := p.lastSeen.Get(context.Background())
	if err != nil {
		log.Printf("Retrieving lastSeen failed: %v", err)
		return
	}

	var lastSeen map[string]time.Time
	err = doc.DataTo(&lastSeen)
	if err != nil {
		log.Printf("Converting lastSeen failed: %v", err)
		return
	}

	deletes := []firestore.Update{}
	for uuid, t := range lastSeen {
		if t.Add(2 * time.Minute).After(time.Now().UTC()) {
			continue
		}
		deletes = append(deletes, firestore.Update{
			Path:  uuid,
			Value: firestore.Delete,
		})
	}

	if len(deletes) == 0 {
		return
	}

	println("LastSeen")

	ctx := context.Background()
	err = p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		err := tx.Update(p.store, deletes)
		if err != nil {
			return err
		}
		return tx.Update(p.lastSeen, deletes)
	})
	if err != nil {
		log.Printf("Deleting timed out failed: %v", err)
		return
	}

	func() {
		p.deletedMutex.Lock()
		defer p.deletedMutex.Unlock()
		for _, d := range deletes {
			p.deleted[d.Path] = true
		}
	}()

	println("SEENANDGONE")
}

func (p *playerServer) Add(ctx context.Context, req *pb.PlayerAddRequest) (*pb.PlayerAddReply, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Generating UUID failed: %v", err)
		return nil, status.Error(codes.Internal, "Could no generate UUID for Player")
	}
	player := &pb.Player{
		Name: req.Name,
		Uuid: uuid.String(),
	}

	bp, err := proto.Marshal(player)
	if err != nil {
		log.Printf("Marshal failed: %v", err)
		return nil, status.Error(codes.Internal, "Serializing Player to backend failed")
	}

	var stored map[string][]byte

	err = p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(p.store)
		if err != nil {
			return fmt.Errorf("Getting Snapshot %v failed: %v", p.store, err)
		}

		err = doc.DataTo(&stored)
		if err != nil {
			return fmt.Errorf("Loading from Snapshot %v failed: %v", doc, err)
		}

		err = tx.Set(p.store, map[string]interface{}{
			uuid.String(): bp,
		}, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("Setting: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Storing Player failed: %v", err)
		return nil, status.Error(codes.Internal, "Storing Player in backend failed")
	}

	others, err := p.unmarshallPlayersFromMap(stored)
	if err != nil {
		log.Printf("Unmarshalling Players failed: %v", err)
		return nil, status.Error(codes.Internal, "Unmarshalling Players failed")
	}

	p.updateLastSeen(uuid.String(), time.Now().UTC())

	return &pb.PlayerAddReply{
		Player: player,
		Others: others,
	}, nil
}

func (p *playerServer) unmarshallPlayersFromMap(pm map[string][]byte) ([]*pb.Player, error) {

	all := []*pb.Player{}

	pl := &pb.Player{}
	for _, bp := range pm {
		pl.Reset()
		err := proto.Unmarshal(bp, pl)
		if err != nil {
			return nil, err
		}
		all = append(all, pl)
	}

	return all, nil
}

/*
func (p *playerServer) Stream(req *pb.PlayerStreamRequest, stream pb.Players_StreamServer) error {

	func() {
		p.deletedMutex.Lock()
		defer p.deletedMutex.Unlock()
		for uuid := range p.deleted {
			stream.Send(&pb.PlayerEvent{
				Event: &pb.PlayerEvent_Removed{
					Removed: &pb.PlayerRemovedEvent{
						Uuid: uuid,
					},
				},
			})
		}
		for k := range p.deleted {
			delete(p.deleted, k)
		}
	}()



		stream.Send(&pb.PlayerEvent{
			Event: &pb.PlayerEvent_Added{
				Added: &pb.PlayerAddedEvent{
					Player: pl,
				},
			},
		})
	}
	println("DONE")
	return nil
}
*/
func (p *playerServer) updateLastSeen(uuid string, t time.Time) {
	ctx := context.Background()
	_, err := p.lastSeen.Set(ctx, map[string]interface{}{
		uuid: t,
	}, firestore.MergeAll)
	if err != nil {
		log.Printf("Updating last seen for %v failed: %v", uuid, err)
	}
}

func (p *playerServer) Touch(ctx context.Context, req *pb.PlayerTouchRequest) (*pb.PlayerTouchReply, error) {
	p.updateLastSeen(req.Uuid, time.Now().UTC())
	return &pb.PlayerTouchReply{}, nil
}

func (p *playerServer) Remove(ctx context.Context, req *pb.PlayerRemoveRequest) (*pb.PlayerRemoveReply, error) {
	p.updateLastSeen(req.Uuid, time.Unix(0, 0))
	return &pb.PlayerRemoveReply{}, nil
}
