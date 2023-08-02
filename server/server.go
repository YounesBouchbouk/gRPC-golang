package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) StreamClientGreet(stream pb.GreetService_StreamClientGreetServer) error {
	fmt.Printf("StreamClientGreet function  was invoked streming request \n")

	result := ""
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			response := &pb.StreamClientResponse{
				Result: result,
			}
			return stream.SendAndClose(response)
		}

		if err != nil {
			log.Fatalf("error getting streams %v", err)
		}

		firstname := req.GetResult().GetFirstname()

		result += "Hello" + firstname + "!"

	}

}

func (*server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {

	fmt.Printf("GreetEveryone from server function  was invoked streming request \n")

	for {
		req, err := stream.Recv()

		if err == io.EOF {

			return nil
		}

		if err != nil {
			log.Fatalf("error getting streams %v", err)
			return nil
		}

		firstname := req.GetGreeting().Firstname

		stream.Send(&pb.GreetEveyoneResponse{
			Result: "hello " + firstname,
		})

	}

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
