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
	ListSubmissionByUserID(userID uuid.UUID) ([]entity.Submission, error)
	FindSubmissionByID(id uuid.UUID) (*entity.Submission, error)
	UpdateSubmission(submission *entity.Submission) (*entity.Submission, error)
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

func (r *submissionRepository) ListSubmissionByUserID(userID uuid.UUID) ([]entity.Submission, error) {
	var submissions []entity.Submission

	if err := r.db.Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
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

func (r *submissionRepository) FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error) {
	var submission entity.Submission

	if err := r.db.Where("transaction_id = ?", id).First(&submission).Error; err != nil {
		return &submission, err
	}

	return &submission, nil
}
