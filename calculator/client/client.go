package main

import (
	"fmt"
	"log"
	"context"

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

	CalculateSum(c)
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