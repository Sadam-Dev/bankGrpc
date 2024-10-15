package main

import (
	"bankGrpc/api/proto"
	"bankGrpc/internal/repository"
	"bankGrpc/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Инициализация репозиториев
	userRepo := repository.NewRepository()
	transactionRepo := repository.NewTransactionRepo()

	// Инициализация сервиса
	bankService := service.NewBankService(userRepo, transactionRepo)

	// Настройка gRPC сервера
	grpcServer := grpc.NewServer()
	proto.RegisterBankServiceServer(grpcServer, bankService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server is running on port :50051")
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
