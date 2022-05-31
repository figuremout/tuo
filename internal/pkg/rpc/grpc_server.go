package rpc

import (
	context "context"
	"log"
	"net"

	pb "github.com/githubzjm/tuo/internal/pkg/rpc/hello"
	grpc "google.golang.org/grpc"
)

// implement
type server struct {
	pb.UnimplementedGreeterServer
}

// implement
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "hello" + in.GetName()}, nil
}

func StartGRPCServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
