package service

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	// "github.com/google/uuid"
)

type EventService interface {
    CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
    CreatePricing(pricing *entity.Pricing) error
}

type eventService struct {
    eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
    return &eventService{eventRepo}
}

func (s *eventService) CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
    return s.eventRepo.CreateEvent(event, pricings)
}
func (s *eventService) CreatePricing(pricing *entity.Pricing) error {
    return s.eventRepo.CreatePricing(pricing)
}
