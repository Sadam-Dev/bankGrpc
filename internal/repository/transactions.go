package repository

import (
	"errors"
	"sync"
	"time"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Timestamp time.Time
}

type TransactionRepo struct {
	transactions []Transaction
	mu           sync.Mutex
}

func NewTransactionRepo() *TransactionRepo {
	return &TransactionRepo{
		transactions: []Transaction{},
	}
}

func (r *Repository) Transfer(sender, receiver string, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	senderUser, senderExists := r.users[sender]
	receiverUser, receiverExists := r.users[receiver]

	if !senderExists || !receiverExists {
		return errors.New("user not found")
	}

	if senderUser.Balance < amount {
		return errors.New("insufficient balance")
	}

	senderUser.Balance -= amount
	receiverUser.Balance += amount

	r.transactions = append(r.transactions, Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now(),
	})

	return nil
}

func (r *Repository) GetTransactions() ([]Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.transactions, nil
}
