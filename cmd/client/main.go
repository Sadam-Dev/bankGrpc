package main

import (
	"bankGrpc/api/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewBankServiceClient(conn)

	// Пример регистрации пользователя
	registerResp, err := client.Register(context.Background(), &proto.RegisterRequest{
		Username: "john",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Register response: %v", registerResp.Message)

	// Пример транзакции
	_, err = client.MakeTransaction(context.Background(), &proto.TransactionRequest{
		Sender:   "john",
		Receiver: "alice",
		Amount:   50,
	})
	if err != nil {
		log.Fatalf("could not make transaction: %v", err)
	}
	log.Println("Transaction successful")

	// Пример получения всех пользователей
	usersResp, err := client.GetAllUsers(context.Background(), &proto.GetUsersRequest{})
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}
	log.Printf("Users: %v", usersResp.Users)

	// Пример получения всех транзакций
	transactionsResp, err := client.GetAllTransactions(context.Background(), &proto.GetTransactionsRequest{})
	if err != nil {
		log.Fatalf("could not get transactions: %v", err)
	}
	log.Printf("Transactions: %v", transactionsResp.Transactions)
}
