package service

import (
	"errors"
	"fmt"
	"html"
	"mime/multipart"

	"github.com/Giafn/Depublic/configs"
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/upload"
	"github.com/google/uuid"
)

type submissionService struct {
	submissionRepo repository.SubmissionRepository
	cfg            *configs.Config
}

type SubmissionService interface {
	CreateSubmission(submission *entity.Submission) (*entity.Submission, error)
	ListSubmission(userId uuid.UUID) ([]entity.Submission, error)
	FindSubmissionByID(id uuid.UUID) (*entity.Submission, error)
	UpdateSubmission(submission *entity.Submission) (*entity.Submission, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindSubmissionByTransactionID(id uuid.UUID) (*entity.Submission, error)
	UploadSubmission(file *multipart.FileHeader) (string, error)
	SendEmailSubmission(status string, submission *entity.Submission) error
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository, cfg *configs.Config) SubmissionService {
	return &submissionService{submissionRepo: submissionRepo, cfg: cfg}
}

func (s *submissionService) CreateSubmission(submission *entity.Submission) (*entity.Submission, error) {
	return s.submissionRepo.CreateSubmission(submission)
}

func (s *submissionService) ListSubmission(userId uuid.UUID) ([]entity.Submission, error) {
	user, err := s.submissionRepo.FindUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user.Role == "Admin" {
		return s.submissionRepo.ListSubmission()
	}
	submissions, err := s.submissionRepo.ListSubmissionByUserID(userId)
	if err != nil {
		return nil, err
	}

	for i := range submissions {
		filename := fmt.Sprintf("http://%s:%s/app/api/v1/file/%s", s.cfg.Host, s.cfg.Port, submissions[i].Filename)

		submissions[i].Filename = filename
	}

	return submissions, nil
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

func (s *submissionService) SendEmailSubmission(status string, submission *entity.Submission) error {
	transaction, err := s.FindTransactionByID(submission.TransactionID)
	if err != nil {
		return err
	}

	user, err := s.FindUserByID(transaction.UserID)
	if err != nil {
		return err
	}
	html := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>Your submission with the name %s has been %s</p>
		<p>Thank you</p>
	`, user.Profiles.FullName, html.EscapeString(submission.Name), status)

	if status == "accepted" {
		html += fmt.Sprintf(`
			<p>Please pay your Ticket <a href="%s">here</a></p>`,
			transaction.PaymentURL,
		)
	}

	ScheduleEmails(
		user.Email,
		"Submission Status",
		html,
	)

	return nil
}
