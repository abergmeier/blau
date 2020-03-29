package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	pb "github.com/abergmeier/blau/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
)

func init() {

	lis = bufconn.Listen(bufSize)
}

func TestMain(m *testing.M) {

	emulatorString := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorString == "" {
		panic("Tests need to have Firestore emulator running and FIRESTORE_EMULATOR_HOST set")
	}

	os.Exit(m.Run())
}

type testServer struct {
	client      *firestore.Client
	s           *grpc.Server
	lastSeenDoc *firestore.DocumentRef
	storeDoc    *firestore.DocumentRef
	ps          *playerServer
}

func ensureEmptyDoc(client *firestore.Client, docPath string) (*firestore.DocumentRef, error) {
	doc := client.Doc(docPath)

	_, _ = doc.Delete(context.Background())

	_, err := doc.Create(context.Background(), map[string]interface{}{})
	if err == nil {
		return doc, nil
	}

	// There may be an error we care about

	s, ok := status.FromError(err)
	if !ok || s.Code() != codes.AlreadyExists {
		return nil, err
	}

	return doc, nil
}

func newTestServer(projectId string) (*testServer, error) {

	client, err := firestore.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}

	lastSeenDoc, err := ensureEmptyDoc(client, lastSeenDocPath)
	if err != nil {
		client.Close()
		return nil, err
	}

	storeDoc, err := ensureEmptyDoc(client, storeDocPath)
	if err != nil {
		client.Close()
		return nil, err
	}

	ps, err := NewPlayerServer(client, 3*time.Second)
	if err != nil {
		client.Close()
		return nil, err
	}
	s := grpc.NewServer()
	pb.RegisterPlayersServer(s, ps)
	go func() {
		if err := s.Serve(lis); err != nil {
			st, _ := status.FromError(err)
			log.Fatalf("Server exited with error: %v/%v", err, st.Code())
		}
	}()
	return &testServer{
		client:      client,
		s:           s,
		lastSeenDoc: lastSeenDoc,
		storeDoc:    storeDoc,
		ps:          ps,
	}, nil
}

func (ts *testServer) Close() {
	println("GRACEFUL")
	ts.s.GracefulStop()
	println("PS")
	ts.ps.Close()
	_, err := ts.lastSeenDoc.Delete(context.Background())
	if err != nil {
		fmt.Printf("Deleting %v failed: %v", ts.lastSeenDoc, err)
	}
	_, err = ts.storeDoc.Delete(context.Background())
	if err != nil {
		fmt.Printf("Deleting %v failed: %v", ts.storeDoc, err)
	}
	ts.client.Close()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestPlayerAdd(t *testing.T) {

	t.Parallel()

	testForEmulator(t)

	ts, err := newTestServer("emulator-project-id-add")
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlayersClient(conn)
	resp, err := client.Add(ctx, &pb.PlayerAddRequest{
		Name: "Foo",
	})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if len(resp.Others) != 0 {
		t.Errorf("Unexpected other players: %v", resp.Others)
	}

	expected := &pb.Player{
		Name: "Foo",
		Uuid: resp.Player.Uuid,
	}
	if !reflect.DeepEqual(resp.Player, expected) {
		t.Errorf("Response was %v. Expected %v.", resp.Player, expected)
	}
}

func TestPlayerRemove(t *testing.T) {

	t.Parallel()

	testForEmulator(t)

	ts, err := newTestServer("emulator-project-id-remove")
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlayersClient(conn)
	resp, err := client.Add(ctx, &pb.PlayerAddRequest{
		Name: "Foo",
	})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	_, err = client.Remove(ctx, &pb.PlayerRemoveRequest{
		Uuid: resp.Player.Uuid,
	})
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	resp, err = client.Add(ctx, &pb.PlayerAddRequest{
		Name: "Bar",
	})
	if err != nil {
		t.Fatalf("Second Add failed: %v", err)
	}

	if len(resp.Others) != 0 {
		t.Errorf("Unexpected Others: %v", resp.Others)
	}
}

func testForEmulator(t *testing.T) {
	// TODO: implement
}
