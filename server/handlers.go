package main

import (
	"bankGrpc/database"
	"bankGrpc/models"
	"bankGrpc/utils"
	"context"
	"fmt"
	"time"

	pb "bankGrpc/proto"
)

type BankServer struct {
	pb.UnimplementedBankServiceServer
}

func (s *BankServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.GenerateHash(req.Password),
		Balance:  float64(100),
	}

	// Сохраняем пользователя в базу данных
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Преобразуем ID в строку для ответа
	return &pb.SignUpResponse{Id: uint64(user.ID), Message: "User created successfully!"}, nil
}

func (s *BankServer) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.GetAllUsersResponse, error) {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	// Преобразуем пользователей в тип protobuf
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:      uint64(user.ID),
			Name:    user.Name,
			Email:   user.Email,
			Balance: float32(user.Balance),
		})
	}

	return &pb.GetAllUsersResponse{Users: pbUsers}, nil
}

func (s *BankServer) CreateTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	var fromUser, toUser models.User

	// Ищем пользователей по ID
	if err := database.DB.First(&fromUser, req.FromUserId).Error; err != nil {
		return nil, fmt.Errorf("sender not found: %v", err)
	}
	if err := database.DB.First(&toUser, req.ToUserId).Error; err != nil {
		return nil, fmt.Errorf("recipient not found: %v", err)
	}

	// Логируем балансы перед транзакцией
	fmt.Printf("Before transaction: FromUserID: %d, Balance: %.2f, ToUserID: %d, Balance: %.2f\n",
		fromUser.ID, fromUser.Balance, toUser.ID, toUser.Balance)

	// Проверяем баланс
	if fromUser.Balance < float64(req.Amount) {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Обновляем балансы
	fromUser.Balance -= float64(req.Amount)
	toUser.Balance += float64(req.Amount)

	// Логируем балансы после обновления
	fmt.Printf("After transaction: FromUserID: %d, New Balance: %.2f, ToUserID: %d, New Balance: %.2f\n",
		fromUser.ID, fromUser.Balance, toUser.ID, toUser.Balance)

	// Сохраняем изменения в базу данных
	if err := database.DB.Save(&fromUser).Error; err != nil {
		return nil, fmt.Errorf("failed to update sender balance: %v", err)
	}
	if err := database.DB.Save(&toUser).Error; err != nil {
		return nil, fmt.Errorf("failed to update recipient balance: %v", err)
	}

	// Логируем факт сохранения балансов
	fmt.Println("Balances successfully updated")

	// Создаем транзакцию
	transaction := &models.Transaction{
		FromUserID: fromUser.ID,
		ToUserID:   toUser.ID,
		Amount:     float64(req.Amount),
		Timestamp:  time.Now(),
	}

	// Сохраняем транзакцию в базу данных
	if err := database.DB.Create(transaction).Error; err != nil {
		return nil, err
	}

	// Логируем факт успешной транзакции
	fmt.Println("Transaction successfully created")

	return &pb.TransactionResponse{Id: uint64(transaction.ID), Message: "Transaction completed successfully!"}, nil
}

func (s *BankServer) GetAllTransactions(ctx context.Context, req *pb.Empty) (*pb.GetAllTransactionsResponse, error) {
	var transactions []models.Transaction

	// Получаем все транзакции из базы данных
	if err := database.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Преобразуем транзакции в тип protobuf
	var pbTransactions []*pb.Transaction
	for _, tx := range transactions {
		pbTransactions = append(pbTransactions, &pb.Transaction{
			Id:         uint64(tx.ID),
			FromUserId: uint64(tx.FromUserID),
			ToUserId:   uint64(tx.ToUserID),
			Amount:     float32(tx.Amount),
			Timestamp:  tx.Timestamp.Format(time.RFC3339),
		})
	}

	return &pb.GetAllTransactionsResponse{Transactions: pbTransactions}, nil
}
