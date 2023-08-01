package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/YounesBouchbouk/gRPC-training/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("Greet function  was invoked %v", req)
	return &pb.GreetResponse{
		Result: &pb.Greeting{
			Firstname: req.GetResult().GetFirstname(),
			Lastname:  req.GetResult().GetLastname(),
		},
	}, nil
}

func (*server) StreamServerGreet(req *pb.StreamServerRequest, stream pb.GreetService_StreamServerGreetServer) error {
	fmt.Printf("StreamServerGreet function  was invoked %v", req)

	for i := 0; i < 10; i++ {
		result := &pb.Greeting{
			Firstname: "younes",
			Lastname:  "bouchbouk" + strconv.Itoa(i),
		}
		stream.Send(&pb.StreamServerResponse{
			Result: result,
		})
		time.Sleep(1000 * time.Millisecond)
	}

	return nil

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
