package binder

import "github.com/google/uuid"

type TransactionCreateRequest struct {
	EventID        uuid.UUID `json:"event_id" validate:"required"`
	TicketQuantity int       `json:"ticket_quantity" validate:"required"`
	Tickets        []Ticket  `json:"tickets" validate:"required"`
}

type Ticket struct {
	BuyerName string    `json:"buyer_name"`
	PricingId uuid.UUID `json:"pricing_id" validate:"required"`
}

type TransactionUpdateRequest struct {
	EventID        uuid.UUID `json:"event_id" validate:"required"`
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	TicketQuantity int       `json:"ticket_quantity" validate:"required"`
	TotalAmount    int       `json:"total_amount" validate:"required"`
	IsPaid         bool      `json:"is_paid" validate:"required"`
}

type MidtransWebhookRequest struct {
	VaNumbers []struct {
		VaNumber string `json:"va_number"`
		Bank     string `json:"bank"`
	} `json:"va_numbers"`
	TransactionTime   string        `json:"transaction_time"`
	TransactionStatus string        `json:"transaction_status"`
	TransactionID     string        `json:"transaction_id"`
	StatusMessage     string        `json:"status_message"`
	StatusCode        string        `json:"status_code"`
	SignatureKey      string        `json:"signature_key"`
	SettlementTime    string        `json:"settlement_time"`
	PaymentType       string        `json:"payment_type"`
	PaymentAmounts    []interface{} `json:"payment_amounts"`
	OrderID           string        `json:"order_id"`
	MerchantID        string        `json:"merchant_id"`
	GrossAmount       string        `json:"gross_amount"`
	FraudStatus       string        `json:"fraud_status"`
	ExpiryTime        string        `json:"expiry_time"`
	Currency          string        `json:"currency"`
}
