package repository

import (
	"errors"
	"fmt"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction, tickets []entity.Ticket, eventID uuid.UUID, userID uuid.UUID) (*entity.Transaction, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)
	DeleteTransaction(id uuid.UUID) error
	GetPricingByEventID(eventID uuid.UUID, pricingID uuid.UUID) (*entity.Pricing, error)
	GetUsersById(id uuid.UUID) (*entity.User, error)
	GetEventByID(id uuid.UUID) (*entity.Event, error)
	FindTransactionsByUserId(userID uuid.UUID) ([]entity.Transaction, error)
	FindTicketByTransactionID(transactionID uuid.UUID) ([]entity.Ticket, error)
	GetPricingById(pricingID uuid.UUID) (*entity.Pricing, error)
	UpdatePricingRemaining(pricingId uuid.UUID, remaining int) (*entity.Pricing, error)
	FindUnpaidTransactionByUserID(userID uuid.UUID, eventID uuid.UUID) (transId uuid.UUID, err error)
	GetSubmissionByTransactionID(transactionID uuid.UUID) (*entity.Submission, error)
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

func (r *transactionRepository) GetPricingByEventID(eventID uuid.UUID, pricingID uuid.UUID) (*entity.Pricing, error) {
	var pricing entity.Pricing
	if err := r.db.Where("event_id = ?", eventID).
		Where("pricing_id = ?", pricingID).
		First(&pricing).Error; err != nil {
		return nil, errors.New("pricing not found")
	}
	return &pricing, nil
}

func (r *transactionRepository) GetUsersById(id uuid.UUID) (*entity.User, error) {
	user := new(entity.User)
	if err := r.db.Preload("Profiles").
		First(&user, "user_id = ?", id).
		Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetEventByID
func (r *transactionRepository) GetEventByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *transactionRepository) FindTransactionsByUserId(userID uuid.UUID) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Find(&transactions, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindTicketByTransactionID(transactionID uuid.UUID) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	if err := r.db.Find(&tickets, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *transactionRepository) GetPricingById(pricingID uuid.UUID) (*entity.Pricing, error) {
	var pricing entity.Pricing
	if err := r.db.First(&pricing, "pricing_id = ?", pricingID).Error; err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *transactionRepository) UpdatePricingRemaining(pricingId uuid.UUID, remaining int) (*entity.Pricing, error) {
	pricing := new(entity.Pricing)
	if err := r.db.Model(pricing).Where("pricing_id = ?", pricingId).Update("remaining", remaining).Error; err != nil {
		return nil, err
	}
	return pricing, nil
}

func (r *transactionRepository) FindUnpaidTransactionByUserID(userID uuid.UUID, eventID uuid.UUID) (transId uuid.UUID, err error) {
	var transaction entity.Transaction
	transaction.ID = uuid.Nil
	if err := r.db.Where("user_id = ?", userID).
		Where("event_id = ?", eventID).
		Where("is_paid = ?", false).
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
