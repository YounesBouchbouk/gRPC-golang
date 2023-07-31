package main

import (
	"context"
	"fmt"
	"log"

	"github.com/YounesBouchbouk/gRPC-training/pb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Client Word")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("some thing went wrong in connection %v", err)
	}
	// when we done with the connextion we close it , (in the end of the code )
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	req := &pb.GreetRequest{
		Result: &pb.Greeting{
			Firstname: "Younes",
			Lastname:  "Bouchbouk",
		},
	}
	fmt.Printf("Ceated client %v", c)
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error in response  %v", err)
	}

	fmt.Println("hello this come from server ", res.Result)

}
