package rpc

import (
	context "context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
	pb "github.com/githubzjm/tuo/internal/pkg/rpc/hello"
	nm "github.com/githubzjm/tuo/internal/server/nodes"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hello() {
	// init grpc client
	addr := "localhost:55555"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "zjm"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

// pseudo code
func Control(nodeID uint, port string) error {
	// Get IP
	node, err := nm.QueryNode(&models.Node{ID: nodeID})
	if err != nil {
		return err
	}

	// Connect gRPC server
	addr := fmt.Sprintf("%s:%s", node.Host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	c := pb.NewGreeterClient(conn)

	// Call remote func
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "zjm"})
	if err != nil {
		return fmt.Errorf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	return nil
}
