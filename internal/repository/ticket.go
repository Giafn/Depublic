package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
	FindTicketByID(id uuid.UUID) (*entity.Ticket, error)
	UpdateTicket(ticket *entity.Ticket) (*entity.Ticket, error)
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