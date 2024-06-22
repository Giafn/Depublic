package binder

import "github.com/google/uuid"

type CreateSubmission struct {
	TransactionID string `form:"transaction_id" validate:"required"`
	File          string `form:"file"`
}

type ListSubmissionResponse struct {
	ID            uuid.UUID `json:"id"`
	EventID       uuid.UUID `json:"event_id"`
	UserID        uuid.UUID `json:"user_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	Filename      string    `json:"filename"`
}
