package service

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository) TransactionService {
	return &transactionService{transactionRepository: transactionRepository}
}

func (s *transactionService) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	newTransaction, err := s.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}
	return newTransaction, nil
}

func (s *transactionService) FindTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	transaction, err := s.transactionRepository.FindTransactionByID(id)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	existingTransaction, err := s.transactionRepository.FindTransactionByID(transaction.ID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	existingTransaction.EventID = transaction.EventID
	existingTransaction.UserID = transaction.UserID
	existingTransaction.TicketQuantity = transaction.TicketQuantity
	existingTransaction.TotalAmount = transaction.TotalAmount
	existingTransaction.IsPaid = transaction.IsPaid

	updatedTransaction, err := s.transactionRepository.UpdateTransaction(existingTransaction)
	if err != nil {
		return nil, err
	}
	return updatedTransaction, nil
}

func (s *transactionService) FindAllTransactions() ([]entity.Transaction, error) {
	transactions, err := s.transactionRepository.FindAllTransactions()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
