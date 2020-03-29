package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	pb "github.com/abergmeier/blau/api/pb"
	"github.com/abergmeier/blau/pkg/server"
	"google.golang.org/grpc"
)

var (
	projectID = flag.String("project-id", os.Getenv("PROJECT_ID"), "GCP Project ID to manage Firestore")
	port      = flag.Int("port", 0, "Port to listen on connections")
)

func main() {
	flag.Parse()

	if *port == 0 {
		portEnv = os.Getenv("PORT")
		if portEnv != "" {
			*port, err := strconv.Atoi(portEnv)
			if err != nil {
				log.Fatalf("Parsing PORT failed: %v", portEnv)
			}
		}
	}
	

	client, err := firestore.NewClient(context.Background(), *projectID)
	if err != nil {
		log.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ps, err := server.NewPlayerServer(client, 3*time.Second)
	if err != nil {
		log.Fatalf("NewPlayerServer failed: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterPlayersServer(grpcServer, ps)
	grpcServer.Serve(lis)
}
