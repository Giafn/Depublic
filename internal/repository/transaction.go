package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
}

type transactionRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewTransactionRepository(db *gorm.DB, cacheable cache.Cacheable) TransactionRepository {
	return &transactionRepository{db: db, cacheable: cacheable}
}

func (r *transactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	transaction.Auditable = entity.NewAuditable()
	if err := r.db.Create(&transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	transaction.Auditable = entity.UpdateAuditable()
	if err := r.db.Save(&transaction).Error; err != nil {
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
