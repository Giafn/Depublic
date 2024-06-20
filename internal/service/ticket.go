package service

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type TicketService interface {
	CreateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	FindTicketByID(id uuid.UUID) (*entity.Ticket, error)
	UpdateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	ValidateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	FindTicketByBookingNumber(bookingNumber string) (*entity.Ticket, error)
	DeleteTicketById(id uuid.UUID) error
	DeleteTicketByBookingNumber(bookingNumber string) error
	FindAllTickets() ([]entity.Ticket, error)
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepository: ticketRepository,
	}
}

func (s *ticketService) CreateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
    // Implementation logic here, for example:
    // return s.ticketRepository.Save(ticket)

	newTicket, err := s.ticketRepository.CreateTicket(ticket)

	if err != nil {
		return nil, err
	}

	return newTicket, err
}

func (s *ticketService) FindTicketByID(id uuid.UUID) (*entity.Ticket, error) {
	return s.ticketRepository.FindTicketByID(id)
}

func (s *ticketService) UpdateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
	return s.ticketRepository.UpdateTicket(ticket)
}

func (s *ticketService) ValidateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
	validatedTicket, err := s.ticketRepository.ValidateTicket(ticket)
    if err != nil {
        return nil, err
    }
    return validatedTicket, nil
}

func (s *ticketService) FindTicketByBookingNumber(bookingNumber string) (*entity.Ticket, error) {
	return s.ticketRepository.FindTicketByBookingNumber(bookingNumber)
}

func (s *ticketService) DeleteTicketById(id uuid.UUID) error {
	return s.ticketRepository.DeleteTicketById(id)
}

func (s *ticketService) DeleteTicketByBookingNumber(bookingNumber string) error {
	return s.ticketRepository.DeleteTicketByBookingNumber(bookingNumber)
}

func (s *ticketService) FindAllTickets() ([]entity.Ticket, error) {
	tickets, err := s.ticketRepository.FindAllTickets()
	if err != nil {
		return nil, err
	}
	return tickets, nil
}