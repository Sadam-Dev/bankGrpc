package main

import (
	"bankGrpc/database"
	pb "bankGrpc/proto"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	database.InitDB()

	server := grpc.NewServer()
	pb.RegisterBankServiceServer(server, &BankServer{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port :50051")
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
