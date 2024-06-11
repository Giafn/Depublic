package entity

import (
    "github.com/google/uuid"
)

type Transaction struct {
    ID            uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
    EventID       uuid.UUID  `json:"event_id"`
    UserID        uuid.UUID  `json:"user_id"`
    TicketQuantity int      `json:"ticket_quantity"`
    TotalAmount   int       `json:"total_amount"`
    IsPaid        bool      `json:"is_paid" gorm:"default:false"`
    Auditable     Auditable `json:"auditable"`
}


func NewTransaction(eventID, userID uuid.UUID, ticketQuantity, totalAmount int, isPaid bool) *Transaction {
    return &Transaction{
        ID:            uuid.New(),
        EventID:       eventID,
        UserID:        userID,
        TicketQuantity: ticketQuantity,
        TotalAmount:   totalAmount,
        IsPaid:        isPaid,
        Auditable:     NewAuditable(),
    }
}

func UpdateTransaction(id, eventID, userID uuid.UUID, ticketQuantity, totalAmount int, isPaid bool) *Transaction {
    return &Transaction{
        ID:            id,
        EventID:       eventID,
        UserID:        userID,
        TicketQuantity: ticketQuantity,
        TotalAmount:   totalAmount,
        IsPaid:        isPaid,
        Auditable:     UpdateAuditable(),
    }
}
