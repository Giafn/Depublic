package service

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
	CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindPricingByEventID(id uuid.UUID) ([]entity.Pricing, error)
	UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
	UpdateEvent(event *entity.Event) (*entity.Event, error)
	UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	DeleteEvent(eventID uuid.UUID) (bool, error)
	DeletePricing(pricingID uuid.UUID) (bool, error)
}

type eventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo}
}

func (s *eventService) CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
	return s.eventRepository.CreateEvent(event, pricings)
}
func (s *eventService) CreatePricing(pricing *entity.Pricing)  (*entity.Pricing, error) {
	return s.eventRepository.CreatePricing(pricing)
}

func (s *eventService) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	event, err := s.eventRepository.FindEventByID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) FindPricingByEventID(id uuid.UUID) ([]entity.Pricing, error) {
	event, err := s.eventRepository.FindPricingByEventID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
	return s.eventRepository.UpdateEventWithPricing(event, pricings)
}

func (s *eventService) UpdateEvent(event *entity.Event) (*entity.Event, error) {
	return s.eventRepository.UpdateEvent(event)
}

func (s *eventService) UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
	return s.eventRepository.UpdatePricing(pricing)
}

func (s *eventService) DeleteEvent(eventID uuid.UUID) (bool, error) {
	event, err := s.eventRepository.FindEventByID(eventID)

	if err != nil {
		return false, err
	}


	return s.eventRepository.DeleteEvent(event)
}

func (s *eventService) DeletePricing(pricingID uuid.UUID) (bool, error) {
	pricing, err := s.eventRepository.FindPricingByID(pricingID)

	if err != nil {
		return false, err
	}


	return s.eventRepository.DeletePricing(pricing)
}