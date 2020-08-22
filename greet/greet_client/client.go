package main

import (
	"context"
	"fmt"
	"grpc/gRPC/greet/greetpb"
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
	//doClientStreaming(c)

	doBiDirectionalStreaming(c)
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

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Started BiDi streaming RPC...")
	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Result: &greetpb.Greeting{
				FirstName: "Sivaprasad",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Result: &greetpb.Greeting{
				FirstName: "Prasad",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Result: &greetpb.Greeting{
				FirstName: "Bhuvan",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Result: &greetpb.Greeting{
				FirstName: "Charitha",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Result: &greetpb.Greeting{
				FirstName: "Teja",
			},
		},
	}
	// We create a stream by invokingg a client
	stream, err := c.GreetEveryOne(context.Background())

	if err != nil {
		log.Fatalf("Error while creating a stream %v", err)
	}
	waitc := make(chan struct{})
	// we send a bunch of messages to the client(ggo routines)
	go func() {
		// function to send bunch of requests
		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()
	// we receive a bunch messages from the server (go routine)
	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving %v\n", err)
				break
			}

			fmt.Printf("Received %v\n", res.GetResult())
		}
		close(waitc)
	}()
	// block until everythingg is done

	<-waitc
}
