package entity

import (
	"github.com/google/uuid"
)

type Pricing struct {
	PricingId  uuid.UUID `json:"pricing_id" gorm:"type:uuid"`
	EventID   uuid.UUID `json:"event_id" gorm:"type:uuid"`
	Name      string    `json:"name"`
	Quota     int       `json:"quota"`
	Remaining int       `json:"remaining"`
	Fee       int       `json:"fee"`
	Auditable
}

func NewPricing(eventID uuid.UUID, name string, fee int, quota, remaining int) *Pricing {
	return &Pricing{
		PricingId:        uuid.New(),
		EventID:   eventID,
		Name:      name,
		Quota:     quota,
		Remaining: remaining,
		Fee:       fee,
		Auditable: NewAuditable(),
	}
}

func UpdatePricing(id, eventID uuid.UUID, name string, fee int, quota, remaining int) *Pricing {
	return &Pricing{
		PricingId:        id,
		EventID:   eventID,
		Name:      name,
		Quota:     quota,
		Remaining: remaining,
		Fee:       fee,
		Auditable: NewAuditable(),
	}
}
