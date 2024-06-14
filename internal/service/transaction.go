package service

import (
    "errors"
    "github.com/Giafn/Depublic/internal/entity"
    "github.com/Giafn/Depublic/internal/repository"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type TransactionService interface {
    CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
    FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
    UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
    FindAllTransactions() ([]entity.Transaction, error)
    GetPricingByEventID(eventID uuid.UUID) (*entity.Pricing, error)
	DeleteTransaction(id uuid.UUID) error
}

type transactionService struct {
    transactionRepository repository.TransactionRepository
    db                    *gorm.DB
}

func NewTransactionService(transactionRepository repository.TransactionRepository, db *gorm.DB) TransactionService {
    return &transactionService{transactionRepository: transactionRepository, db: db}
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

func (s *transactionService) GetPricingByEventID(eventID uuid.UUID) (*entity.Pricing, error) {
    var pricing entity.Pricing
    if err := s.db.Where("event_id = ?", eventID).First(&pricing).Error; err != nil {
        return nil, err
    }
    return &pricing, nil
}

func (s *transactionService) DeleteTransaction(id uuid.UUID) error {
    return s.transactionRepository.DeleteTransaction(id)
}