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

// json nya
