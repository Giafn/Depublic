package entity

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/martinlindhe/base36"

	"github.com/google/uuid"
)

type Ticket struct {
	ID            uuid.UUID `json:"id"`
	TransactionID string    `json:"transaction_id"`
	EventID       string    `json:"event_id"`
	Name          string    `json:"name"` // nama pemilik tiket
	BookingNum    string    `json:"bookingNum"`
	IsUsed        bool      `json:"isUsed"`
	PricingID     string    `json:"pricing_id"`
	Auditable
}

func NewTicket(transactionID, eventID, name string, pricingID string) *Ticket {
	return &Ticket{
		ID:            uuid.New(),
		TransactionID: transactionID,
		EventID:       eventID,
		Name:          name,
		BookingNum:    createBookingNumber(transactionID, eventID, name),
		IsUsed:        false,
		PricingID:     pricingID,
		Auditable:     NewAuditable(),
	}
}

func createBookingNumber(transactionID, eventID, visitorName string) string {
	// Get the current timestamp in nanoseconds
	timestamp := time.Now().UnixNano()

	// Create a string to hash by concatenating the input values and timestamp
	input := fmt.Sprintf("%s-%s-%s-%d", transactionID, eventID, visitorName, timestamp)

	// Create a SHA-256 hash of the input string
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	// Use the first 8 bytes of the hash for a shorter unique identifier
	shortHashBytes := hashBytes[:8]

	// Convert the hash bytes to a base36 string
	bookingNumber := base36.EncodeBytes(shortHashBytes)

	return bookingNumber
}

func UpdateTicket(oldTicket Ticket, name string) *Ticket {
	return &Ticket{
		ID:            oldTicket.ID,
		TransactionID: oldTicket.TransactionID,
		EventID:       oldTicket.EventID,
		Name:          name,
		BookingNum:    oldTicket.BookingNum,
		IsUsed:        oldTicket.IsUsed,
		PricingID:     oldTicket.PricingID,
		Auditable:     UpdateAuditable(),
	}
}

func ValidateTicket(oldTicket Ticket) *Ticket {
	return &Ticket{
		ID:            oldTicket.ID,
		TransactionID: oldTicket.TransactionID,
		EventID:       oldTicket.EventID,
		Name:          oldTicket.Name,
		BookingNum:    oldTicket.BookingNum,
		IsUsed:        true,
		PricingID:     oldTicket.PricingID,
		Auditable:     UpdateAuditable(),
	}
}
