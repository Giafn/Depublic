package service

import (
	"errors"
	"mime/multipart"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/upload"
	"github.com/google/uuid"
)

type submissionService struct {
	submissionRepo repository.SubmissionRepository
}

type SubmissionService interface {
	CreateSubmission(submission *entity.Submission) (*entity.Submission, error)
	ListSubmission() ([]entity.Submission, error)
	FindSubmissionByID(id uuid.UUID) (*entity.Submission, error)
	UpdateSubmission(submission *entity.Submission) (*entity.Submission, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error)
	UploadSubmission(file *multipart.FileHeader) (string, error)
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository) SubmissionService {
	return &submissionService{submissionRepo}
}

func (s *submissionService) CreateSubmission(submission *entity.Submission) (*entity.Submission, error) {
	return s.submissionRepo.CreateSubmission(submission)
}

func (s *submissionService) ListSubmission() ([]entity.Submission, error) {
	return s.submissionRepo.ListSubmission()
}

func (s *submissionService) FindSubmissionByID(id uuid.UUID) (*entity.Submission, error) {
	return s.submissionRepo.FindSubmissionByID(id)
}

func (s *submissionService) UpdateSubmission(submission *entity.Submission) (*entity.Submission, error) {
	return s.submissionRepo.UpdateSubmission(submission)
}

func (s *submissionService) FindTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	return s.submissionRepo.FindTransactionByID(id)
}

func (s *submissionService) FindUserByID(id uuid.UUID) (*entity.User, error) {
	return s.submissionRepo.FindUserByID(id)
}

func (s *submissionService) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	return s.submissionRepo.FindEventByID(id)
}

func (s *submissionService) FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error) {
	return s.submissionRepo.FindSubmissionByTransactionID(id)
}

func (s *submissionService) UploadSubmission(file *multipart.FileHeader) (string, error) {
	if file != nil {
		filename, err := upload.UploadFile(file, "submissions")
		if err != nil {
			return "", err
		}
		return filename, nil
	}
	return "", errors.New("file is nil")
}
