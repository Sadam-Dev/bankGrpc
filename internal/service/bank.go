package service

import (
	"bankGrpc/api/proto"
	"bankGrpc/internal/repository"
	"bankGrpc/validation"
	"context"
	"time"
)

type BankService struct {
	userRepo        *repository.Repository
	transactionRepo *repository.TransactionRepo
}

func NewBankService(userRepo *repository.Repository, transactionRepo *repository.TransactionRepo) *BankService {
	return &BankService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *BankService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if err := validation.ValidateUser(req.Username, req.Password); err != nil {
		return nil, err
	}
	err := s.userRepo.CreateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *BankService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.userRepo.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}

func (s *BankService) MakeTransaction(ctx context.Context, req *proto.TransactionRequest) (*proto.TransactionResponse, error) {
	if err := validation.ValidateTransaction(req.Amount); err != nil {
		return nil, err
	}

	err := s.userRepo.Transfer(req.Sender, req.Receiver, req.Amount)
	if err != nil {
		return nil, err
	}
	return &proto.TransactionResponse{Message: "Transaction completed successfully"}, nil
}

func (s *BankService) GetAllUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	users, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.User
	for _, u := range users {
		protoUsers = append(protoUsers, &proto.User{
			Username: u.Username,
			Balance:  u.Balance,
		})
	}

	return &proto.GetUsersResponse{Users: protoUsers}, nil
}

func (s *BankService) GetAllTransactions(ctx context.Context, req *proto.GetTransactionsRequest) (*proto.GetTransactionsResponse, error) {
	transactions, err := s.userRepo.GetTransactions()
	if err != nil {
		return nil, err
	}

	var protoTransactions []*proto.Transaction
	for _, t := range transactions {
		protoTransactions = append(protoTransactions, &proto.Transaction{
			Sender:    t.Sender,
			Receiver:  t.Receiver,
			Amount:    t.Amount,
			Timestamp: t.Timestamp.Format(time.RFC3339),
		})
	}

	return &proto.GetTransactionsResponse{Transactions: protoTransactions}, nil
}
