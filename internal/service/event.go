package service

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
	CreatePricing(pricing *entity.Pricing) error
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindPricingByEventID(id uuid.UUID) ([]entity.Pricing, error)
    GetEvents(filters map[string]interface{}, sort string) ([]entity.Event,error)
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
func (s *eventService) CreatePricing(pricing *entity.Pricing) error {
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

func (s *eventService) GetEvents(filters map[string]interface{}, sort string) ([]entity.Event,error){
    events, err := s.eventRepository.GetEvents(filters,sort)
    if err != nil {
        return nil, err
    }
    return events, nil
}