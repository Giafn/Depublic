package entity

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	EventID        uuid.UUID `json:"event_id"`
	UserID         uuid.UUID `json:"user_id"`
	TicketQuantity int       `json:"ticket_quantity"`
	TotalAmount    int       `json:"total_amount"`
	IsPaid         bool      `json:"is_paid" gorm:"default:false"`
	PaymentURL     string    `json:"payment_url"`
	Auditable
}

type Events struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Organizer   string    `json:"organizer"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	MustUpload  bool      `json:"must_upload_submission"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	District    string    `json:"district"`
	FullAddress string    `json:"full_address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Auditable
}

func NewTransaction(eventID, userID uuid.UUID, ticketQuantity, totalAmount int, isPaid bool) *Transaction {
	return &Transaction{
		ID:             uuid.New(),
		EventID:        eventID,
		UserID:         userID,
		TicketQuantity: ticketQuantity,
		TotalAmount:    totalAmount,
		IsPaid:         isPaid,
		Auditable:      NewAuditable(),
	}
}

func UpdateTransaction(id, eventID, userID uuid.UUID, ticketQuantity, totalAmount int, isPaid bool) *Transaction {
	return &Transaction{
		ID:             id,
		EventID:        eventID,
		UserID:         userID,
		TicketQuantity: ticketQuantity,
		TotalAmount:    totalAmount,
		IsPaid:         isPaid,
		Auditable:      UpdateAuditable(),
	}
}
