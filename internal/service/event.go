package service

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
	CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	FindEventByID(id uuid.UUID) (*entity.Event, error)
	FindPricingByEventID(id uuid.UUID) ([]entity.Pricing, error)
	GetEvents(filters map[string]interface{}, sort string, distance map[string]float64, pagination map[string]int) ([]entity.Event, int, error)
	UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error)
	UpdateEvent(event *entity.Event) (*entity.Event, error)
	UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error)
	DeleteEvent(eventID uuid.UUID) (bool, error)
	DeletePricing(pricingID uuid.UUID) (bool, error)
}

type eventService struct {
	eventRepository   repository.EventRepository
	pricingRepository repository.PricingRepository
}

func NewEventService(eventRepo repository.EventRepository, pricingRepo repository.PricingRepository) EventService {
	return &eventService{eventRepository: eventRepo, pricingRepository: pricingRepo}
}

func (s *eventService) CreateEvent(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
	return s.eventRepository.CreateEvent(event, pricings)
}
func (s *eventService) CreatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
	return s.pricingRepository.CreatePricing(pricing)
}

func (s *eventService) FindEventByID(id uuid.UUID) (*entity.Event, error) {
	event, err := s.eventRepository.FindEventByID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) FindPricingByEventID(id uuid.UUID) ([]entity.Pricing, error) {
	event, err := s.pricingRepository.FindPricingByEventID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) GetEvents(filters map[string]interface{}, sort string, distance map[string]float64, pagination map[string]int) ([]entity.Event, int, error) {
	if sort == "terdekat" && distance == nil {
		return nil, 0, errors.New("latitude dan longitude diperlukan untuk sort terdekat")
	}
	events, count, err := s.eventRepository.GetEvents(filters, sort, distance, pagination)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}
func (s *eventService) UpdateEventWithPricing(event *entity.Event, pricings []entity.Pricing) (*entity.Event, error) {
	return s.eventRepository.UpdateEventWithPricing(event, pricings)
}

func (s *eventService) UpdateEvent(event *entity.Event) (*entity.Event, error) {
	return s.eventRepository.UpdateEvent(event)
}

func (s *eventService) UpdatePricing(pricing *entity.Pricing) (*entity.Pricing, error) {
	return s.pricingRepository.UpdatePricing(pricing)
}

func (s *eventService) DeleteEvent(eventID uuid.UUID) (bool, error) {
	event, err := s.eventRepository.FindEventByID(eventID)

	if err != nil {
		return false, err
	}

	return s.eventRepository.DeleteEvent(event)
}

func (s *eventService) DeletePricing(pricingID uuid.UUID) (bool, error) {
	pricing, err := s.pricingRepository.FindPricingByID(pricingID)

	if err != nil {
		return false, err
	}

	return s.pricingRepository.DeletePricing(pricing)
}
