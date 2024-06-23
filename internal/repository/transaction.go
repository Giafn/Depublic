package repository

import (
	"errors"
	"fmt"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// submissionRepository

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction, tickets []entity.Ticket, eventID uuid.UUID, userID uuid.UUID) (*entity.Transaction, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)
	DeleteTransaction(id uuid.UUID) error
	FindTransactionsByUserId(userID uuid.UUID) ([]entity.Transaction, error)
	// FindTicketByTransactionID(transactionID uuid.UUID) ([]entity.Ticket, error)
	FindUnpaidTransactionByUserID(userID uuid.UUID, eventID uuid.UUID) (transId uuid.UUID, err error)
	GetSubmissionByTransactionID(transactionID uuid.UUID) (*entity.Submission, error)
	UpdateTransactionStatus(transactionID uuid.UUID, status string) (*entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(
	transaction *entity.Transaction,
	tickets []entity.Ticket,
	eventID uuid.UUID,
	userID uuid.UUID,
) (*entity.Transaction, error) {
	tx := r.db.Begin()
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range tickets {
		tickets[i].TransactionID = transaction.ID.String()
		if err := tx.Create(&tickets[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return transaction, nil
}

func (r *transactionRepository) FindTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if err := r.db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) FindAllTransactions() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) DeleteTransaction(id uuid.UUID) error {
	return r.db.Delete(&entity.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) FindTransactionsByUserId(userID uuid.UUID) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Find(&transactions, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindUnpaidTransactionByUserID(userID uuid.UUID, eventID uuid.UUID) (transId uuid.UUID, err error) {
	var transaction entity.Transaction
	transaction.ID = uuid.Nil
	if err := r.db.Where("user_id = ?", userID).
		Where("event_id = ?", eventID).
		Where("is_paid = ?", false).
		Where("status = ?", "pending").
		First(&transaction).Error; err != nil {
		return uuid.Nil, nil
	}

	if transaction.ID != uuid.Nil {
		return transaction.ID, fmt.Errorf("transaction with id %s is unpaid", transaction.ID)
	}
	return uuid.Nil, nil
}

func (r *transactionRepository) GetSubmissionByTransactionID(transactionID uuid.UUID) (*entity.Submission, error) {
	var submission entity.Submission
	if err := r.db.First(&submission, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, errors.New("please upload your submission")
	}
	return &submission, nil
}

func (r *transactionRepository) UpdateTransactionStatus(transactionID uuid.UUID, status string) (*entity.Transaction, error) {
	transaction, err := r.FindTransactionByID(transactionID)
	if err != nil {
		return nil, err
	}
	transaction.Status = status
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
