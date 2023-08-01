package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/YounesBouchbouk/gRPC-training/pb"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("some thing went wrong in connection %v", err)
	}
	// when we done with the connextion we close it , (in the end of the code )
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	// unary_grpc(c)
	server_stream_grpc(c)
}

func unary_grpc(c pb.GreetServiceClient) {

	req := &pb.GreetRequest{
		Result: &pb.Greeting{
			Firstname: "Younes",
			Lastname:  "Bouchbouk",
		},
	}
	fmt.Printf("Created client %v", c)
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error in response  %v", err)
	}

	fmt.Println("hello this come from server ", res.Result)
}

func server_stream_grpc(c pb.GreetServiceClient) {

	req := &pb.StreamServerRequest{
		Result: &pb.Greeting{
			Firstname: "younes",
			Lastname:  "bouchbouk",
		},
	}
	resStreem, err := c.StreamServerGreet(context.Background(), req)

	if err != nil {
		log.Fatalf("error with calling StrelServerGreet  %v", err)
	}

	for {
		msg, err := resStreem.Recv()

		if err == io.EOF {
			//t that's mean that the server has ended
			break
		}

		if err != nil {
			log.Fatalf("error while reading strem  %v", err)
		}

		log.Printf("response from GreetManytimes : %v", msg.GetResult())

	}

}
