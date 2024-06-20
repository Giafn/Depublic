package repository

import (
	"errors"

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
	GetEventByID(id uuid.UUID) (*entity.Events, error)
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
	var user entity.User
	if err := r.db.First(&user, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetEventByID
func (r *transactionRepository) GetEventByID(id uuid.UUID) (*entity.Events, error) {
	var event entity.Events
	if err := r.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
