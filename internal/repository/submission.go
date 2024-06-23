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
	ListSubmission(page, limit int) ([]entity.Submission, int, error)
	ListSubmissionByUserID(userID uuid.UUID, page, limit int) ([]entity.Submission, int, error)
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

func (r *submissionRepository) ListSubmission(page, limit int) ([]entity.Submission, int, error) {
	var submissions []entity.Submission

	if err := r.db.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&submissions).Error; err != nil {
		return submissions, 0, err
	}
	var count int64
	if err := r.db.Model(&entity.Submission{}).Count(&count).Error; err != nil {
		return submissions, 0, err
	}
	return submissions, int(count), nil
}

func (r *submissionRepository) ListSubmissionByUserID(userID uuid.UUID, page, limit int) ([]entity.Submission, int, error) {
	var submissions []entity.Submission

	if err := r.db.
		Offset((page-1)*limit).
		Limit(limit).
		Where("user_id = ?", userID).
		Find(&submissions).Error; err != nil {
		return submissions, 0, err
	}

	var count int64
	if err := r.db.Model(&entity.Submission{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return submissions, 0, err
	}

	return submissions, int(count), nil
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
