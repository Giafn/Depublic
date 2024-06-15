package entity

import (
    "github.com/google/uuid"
)

type Pricing struct {
    PricingID uuid.UUID `json:"pricing_id" gorm:"type:uuid;default:uuid_generate_v4()"`
    EventID   uuid.UUID `json:"event_id" gorm:"type:uuid"`
    Name      string    `json:"name"`
    Fee       int       `json:"fee"`
}
