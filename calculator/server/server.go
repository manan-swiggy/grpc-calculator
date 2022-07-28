package main

import (
	"fmt"
	"net"
	"context"
	"log"
	"example.com/calculator/calculatorpb"

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