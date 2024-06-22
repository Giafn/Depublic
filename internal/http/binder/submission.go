package binder

type CreateSubmission struct {
	TransactionID string `form:"transaction_id" validate:"required"`
	File          string `form:"file"`
}
