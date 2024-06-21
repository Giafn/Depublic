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
	IDTransaction string `json:"idTransaction"`
	IDEvent       string `json:"idEvent"`
	Name          string `json:"name"` // nama pemilik tiket
	BookingNum    string `json:"bookingNum"`
	IsUsed        bool   `json:"isUsed"`
	Auditable
}

func NewTicket(idTransaction, idEvent, name string) *Ticket {
	return &Ticket{
		ID: uuid.New(),
		IDTransaction: idTransaction,
		IDEvent: idEvent,
		Name: name,
		BookingNum: createBookingNumber(idTransaction, idEvent, name),
		IsUsed: false,
        Auditable:     NewAuditable(),
	}
}

func createBookingNumber(idTransaction, idEvent, visitorName string) string {
	// Get the current timestamp in nanoseconds
	timestamp := time.Now().UnixNano()

	// Create a string to hash by concatenating the input values and timestamp
	input := fmt.Sprintf("%s-%s-%s-%d", idTransaction, idEvent, visitorName, timestamp)

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
		ID: oldTicket.ID,
		IDTransaction: oldTicket.IDTransaction,
		IDEvent: oldTicket.IDEvent,
		Name: name,
		BookingNum: oldTicket.BookingNum,
		IsUsed: oldTicket.IsUsed,
        Auditable:     UpdateAuditable(),
	}
}

func ValidateTicket(oldTicket Ticket) *Ticket {
	return &Ticket{
		ID: oldTicket.ID,
		IDTransaction: oldTicket.IDTransaction,
		IDEvent: oldTicket.IDEvent,
		Name: oldTicket.Name,
		BookingNum: oldTicket.BookingNum,
		IsUsed: true,
        Auditable:     UpdateAuditable(),

	}
}
