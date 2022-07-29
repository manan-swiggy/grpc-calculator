package main

import (
	"fmt"
	"log"
	"context"
	"io"
	"time"

	"example.com/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main () {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	fmt.Println("starting client")

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// CalculateSum(c)

	// GetSmallerPrimes(c)

	// ComputeAverage(c)

	FindMaximum(c)
}

func CalculateSum (c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Sum grpc ...")

	req := calculatorpb.SumRequest {
		Sum: &calculatorpb.Sum {
			Num1: 12,
			Num2: 3,
		},
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling sum grpc unary: %v", err)
	}

	log.Printf("Response from Sum Unary call: %v", resp.Result)
}

func GetSmallerPrimes(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting server side grpc streaming")

	req := calculatorpb.ReturnPrimesRequest{
		Num: 15,
	}

	respStream, err := c.ReturnSmallerPrimes(context.Background(), &req)
	if err != nil {
		log.Fatal("error while calling GetSmallerPrimes server side streaming: %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("error whilw receiving server stream: %v", err)
		}

		fmt.Println("Response from GetSmallerPrimes server: %v", msg.Result)
	}
}

func ComputeAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting client side streaming over grpc ...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Num: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 44,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 467,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 6,
		},
	}

	for _, req := range requests {
		fmt.Println("\nSending Request.... : ", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}

	fmt.Println("\n****Response From Server : ", resp.GetResult())
}

func FindMaximum (c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting bi-directional streaming by calling FindMaximum over GRPC ...")

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest {
			Num: 12,
		},
		&calculatorpb.FindMaximumRequest {
			Num: 34,
		},
		&calculatorpb.FindMaximumRequest {
			Num: 99,
		},
		&calculatorpb.FindMaximumRequest {
			Num: 1,
		},
	}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming: %v", err)
	}

	waitchan := make(chan struct{})

	go func (requests []*calculatorpb.FindMaximumRequest) {
		for _, req := range requests {
			fmt.Println("\n Sending Request ...", req.Num)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request to FindMaximum service : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func () {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("error receiving reponse from server: %v", err)
			}

			fmt.Println("\n Response from Server: ", resp.GetResult())
		}
	} ()

	<- waitchan
}