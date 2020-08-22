package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//fmt.Printf("Created client : %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sivaprasad",
			LastName:  "Tamatam",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}

	log.Printf("Response from Greet : %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Start to do Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sivaprasad",
			LastName:  "Tamatam",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error in streaming %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading %v", err)
		}

		fmt.Println("Result " + msg.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Started Client streaming RPC...")
	stream, err := c.Longgreet(context.Background())
	requests := []*greetpb.LogGreetRequest{
		&greetpb.LogGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sivaprasad",
			},
		},
		&greetpb.LogGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Prasad",
			},
		},
		&greetpb.LogGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bhuvan",
			},
		},
		&greetpb.LogGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Charitha",
			},
		},
		&greetpb.LogGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Teja",
			},
		},
	}
	if err != nil {
		log.Fatalf("Error while calling LongGreet %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receivong %v", err)
	}
	fmt.Printf("%v", res)
}
