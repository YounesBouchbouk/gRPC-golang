package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// server_stream_grpc(c)
	client_stream_grpc(c)
}

func unary_grpc(c pb.GreetServiceClient) {

	fmt.Println("unary rpc has been invoked ")

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

	fmt.Println("server_stream_grpc has been invoked")

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

func client_stream_grpc(c pb.GreetServiceClient) {

	requestsList := []*pb.StreamClientRequest{
		&pb.StreamClientRequest{
			Result: &pb.Greeting{
				Firstname: "younes",
				Lastname:  "bouche",
			},
		},
		&pb.StreamClientRequest{
			Result: &pb.Greeting{
				Firstname: "med",
				Lastname:  "sb3",
			},
		},
		&pb.StreamClientRequest{
			Result: &pb.Greeting{
				Firstname: "l3ba",
				Lastname:  "j3iba",
			},
		},
	}

	stream, err := c.StreamClientGreet(context.Background())

	if err != nil {
		log.Fatalf("error while connecting to server  %v \n", err)
	}

	for _, req := range requestsList {
		fmt.Printf("sending this data from client   %v \n", req)

		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("cannot get from server  %v \n", err)
	}

	fmt.Printf("result has been rexeived successfully %v \n", resp)

}
