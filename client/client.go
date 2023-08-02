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
	// client_stream_grpc(c)
	bidi_grpc(c)
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

func bidi_grpc(c pb.GreetServiceClient) {

	requestsList := []*pb.GreetEveyoneRequest{
		&pb.GreetEveyoneRequest{
			Greeting: &pb.Greeting{
				Firstname: "younes",
				Lastname:  "bouche",
			},
		},
		&pb.GreetEveyoneRequest{
			Greeting: &pb.Greeting{
				Firstname: "med",
				Lastname:  "sb3",
			},
		},
		&pb.GreetEveyoneRequest{
			Greeting: &pb.Greeting{
				Firstname: "l3ba",
				Lastname:  "j3iba",
			},
		},
	}

	//get the stream object
	stream, err := c.GreetEveryone(context.Background())

	//channel to wait for routine to complelte
	waitChan := make(chan struct{})

	if err != nil {
		log.Fatalf("error while connecting to server  %v \n", err)
	}

	//goroutine function to receive from server

	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				//t that's mean that the server has ended
				break
			}

			if err != nil {
				log.Fatalf("error while reading strem  %v", err)
				break
			}

			log.Printf("response from bidiGrpc : %v", msg.GetResult())

		}
		close(waitChan)
	}()

	//goroutine function to send to server

	go func() {
		for _, req := range requestsList {
			Serr := stream.Send(req)
			time.Sleep(400 * time.Millisecond)
			if Serr != nil {
				log.Fatalf("error while sending strem  %v", err)
				return
			}
		}

	}()

	//close the channel to stop the client
	<-waitChan
}
