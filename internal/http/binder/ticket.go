package binder

import "github.com/Giafn/Depublic/internal/entity"

type TicketCreateRequest struct {
	EventID       string `json:"event_id" validate:"required"`
	TransactionID string `json:"transaction_id" validate:"required"`
	Data          []entity.Person `json:"data" validate:"required"`
}

type TicketUpdateRequest struct {
	ID string `param:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type TicketFindByIdRequest struct {
	ID string `param:"id" validate:"required"`
}

type TicketValidateRequest struct {
	ID string `param:"id" validate:"required"`
}

type TicketFindByBookingNumberRequest struct {
	BookingNumber string `param:"bookingNum" validate:"required"`
}