package main

import (
	"fmt"
	"net"
	"context"
	"log"
	"time"
	"example.com/calculator/calculatorpb"
	"example.com/calculator/helper"
	"io"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum (ctx context.Context, req *calculatorpb.SumRequest) (resp *calculatorpb.SumResponse, err error) {
	fmt.Println("Sum function was invoked to add")

	num1 := req.GetSum().GetNum1()
	num2 := req.GetSum().GetNum2()

	res := num1 + num2
	resp = &calculatorpb.SumResponse{
		Result: res,
	}
	return resp, nil
}

func (*server) ReturnSmallerPrimes(req * calculatorpb.ReturnPrimesRequest, resp calculatorpb.CalculatorService_ReturnSmallerPrimesServer) error {
	fmt.Println("GetSamllerPrimes function invoked for server side streaming")

	primes := helper.Sieve(int(req.Num))

	for _, p := range primes {
		res := &calculatorpb.ReturnPrimesResponse {
			Result: int64(p),
		}

		time.Sleep(1000 * time.Millisecond)
		resp.Send(res)
	}
	return nil
}

func (*server) ComputeAverage (stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage Function is invoked to demonstrate client side streaming")

	var totalSum int64
	var totalCount int64

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: totalSum/totalCount,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
	
		totalCount++
		totalSum += msg.Num
	}
	
}

func main () {
	fmt.Println("vim-go")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}