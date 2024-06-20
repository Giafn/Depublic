package repository

import (
    "gorm.io/gorm"
    "github.com/Giafn/Depublic/internal/entity"
	// "github.com/google/uuid"
)

type EventRepository interface {
    CreateEvent(event *entity.Event) error
    CreatePricing(pricing *entity.Pricing) error
}

type eventRepository struct {
    db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
    return &eventRepository{db}
}

func (r *eventRepository) CreateEvent(event *entity.Event) error {
    return r.db.Create(event).Error
}


func (r *eventRepository) CreatePricing(pricing *entity.Pricing) error {
    return r.db.Create(pricing).Error
}