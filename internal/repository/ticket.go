package repository

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	FindTicketByID(id uuid.UUID) (*entity.Ticket, error)
	UpdateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	ValidateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	FindTicketByBookingNumber(bookingNumber string) (*entity.Ticket, error)
	DeleteTicketById(id uuid.UUID) error
	DeleteTicketByBookingNumber(bookingNumber string) error
	FindAllTickets() ([]entity.Ticket, error)
	FindTicketsByTransactionId(transactionId uuid.UUID) ([]entity.Ticket, error)
}
	
type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (r *ticketRepository) CreateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
	if err := r.db.Create(&ticket).Error; err != nil {
		return ticket, err
	}
	return ticket, nil
}

func (r *ticketRepository) FindTicketByID(id uuid.UUID) (*entity.Ticket, error) {
	ticket := new(entity.Ticket)

	if err := r.db.Where("tickets.id = ?", id).Take(ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func (r *ticketRepository) UpdateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
	fields := make(map[string]interface{})

	if ticket.Name != "" {
		fields["name"] = ticket.Name
	}

	if err := r.db.Model(ticket).Where("id = ?", ticket.ID).Updates(fields).Error; err != nil {
		return ticket, nil
	}

	return ticket, nil
}

var ErrTicketAlreadyValidated = errors.New("ticket is already validated")

func (r *ticketRepository) ValidateTicket(ticket *entity.Ticket) (*entity.Ticket, error) {
	if ticket.IsUsed {
		return ticket, ErrTicketAlreadyValidated
	}
	
	fields := make(map[string]interface{})

	fields["is_used"] = ticket.IsUsed

	if err := r.db.Model(ticket).Where("id = ?", ticket.ID).Updates(fields).Error; err != nil {
		return ticket, nil
	}

	return ticket, nil
}

func (r *ticketRepository) FindTicketByBookingNumber(bookingNumber string) (*entity.Ticket, error) {
	ticket := new(entity.Ticket)

	if err := r.db.Where("tickets.booking_num = ?", bookingNumber).Take(ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func (r *ticketRepository) DeleteTicketById(id uuid.UUID) error {

	ticket := new(entity.Ticket)

	if err := r.db.Where("tickets.id = ?", id).Take(ticket).Error; err != nil {
		return err
	}

	if err := r.db.Where("tickets.id = ?", id).Delete(ticket).Error; err != nil {
		return err
	}

	return nil
}

func (r *ticketRepository) DeleteTicketByBookingNumber(bookingNumber string) error {

	ticket := new(entity.Ticket)

	if err := r.db.Where("tickets.booking_num = ?", bookingNumber).Take(ticket).Error; err != nil {
		return err
	}

	if err := r.db.Where("tickets.booking_num = ?", bookingNumber).Delete(ticket).Error; err != nil {
		return err
	}

	return nil
}

func (r *ticketRepository) FindAllTickets() ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	if err := r.db.Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) FindTicketsByTransactionId(transactionId uuid.UUID) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
    result := r.db.Where("id_transaction = ?", transactionId).Find(&tickets)
    if result.Error != nil {
        return nil, result.Error
    }
    return tickets, nil
}