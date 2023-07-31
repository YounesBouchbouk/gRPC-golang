package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/YounesBouchbouk/gRPC-training/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("Greet Funxtion was invoked %v", req)
	return &pb.GreetResponse{
		Result: &pb.Greeting{
			Firstname: req.GetResult().GetFirstname(),
			Lastname:  req.GetResult().GetLastname(),
		},
	}, nil
}

func main() {

	fmt.Println("Server world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to lister : %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start : %v", err)

	}

}
