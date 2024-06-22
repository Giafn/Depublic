package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type submissionRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

type SubmissionRepository interface {
	CreateSubmission(submission *entity.Submission) (*entity.Submission, error)
	ListSubmission() ([]entity.Submission, error)
	FindSubmissionByID(id uuid.UUID) (*entity.Submission, error)
	UpdateSubmission(submission *entity.Submission) (*entity.Submission, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error)
}

func NewSubmissionRepository(db *gorm.DB, cacheable cache.Cacheable) SubmissionRepository {
	return &submissionRepository{db: db, cacheable: cacheable}
}

func (r *submissionRepository) CreateSubmission(submission *entity.Submission) (*entity.Submission, error) {
	submission.ID = uuid.New()
	if err := r.db.Create(&submission).Error; err != nil {
		return submission, err
	}

	return submission, nil
}

func (r *submissionRepository) ListSubmission() ([]entity.Submission, error) {
	var submissions []entity.Submission

	if err := r.db.Find(&submissions).Error; err != nil {
		return submissions, err
	}

	return submissions, nil
}

func (r *submissionRepository) FindSubmissionByID(id uuid.UUID) (*entity.Submission, error) {
	var submission entity.Submission

	if err := r.db.First(&submission, id).Error; err != nil {
		return &submission, err
	}

	return &submission, nil
}

func (r *submissionRepository) UpdateSubmission(submission *entity.Submission) (*entity.Submission, error) {

	if err := r.db.Save(&submission).Error; err != nil {
		return submission, err
	}

	return submission, nil
}

func (r *submissionRepository) FindTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction

	if err := r.db.First(&transaction, id).Error; err != nil {
		return &transaction, err
	}

	return &transaction, nil
}

func (r *submissionRepository) FindUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User

	if err := r.db.Preload("Profiles").
		First(&user, id).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

func (r *submissionRepository) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event

	if err := r.db.First(&event, id).Error; err != nil {
		return &event, err
	}

	return &event, nil
}

func (r *submissionRepository) FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error) {
	var submission entity.Submission

	if err := r.db.Where("transaction_id = ?", id).First(&submission).Error; err != nil {
		return &submission, err
	}

	return &submission, nil
}
