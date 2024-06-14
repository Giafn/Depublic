package binder

import "github.com/Giafn/Depublic/internal/entity"

type TicketCreateRequest struct {
	IDEvent       string `json:"idEvent"`
	IDTransaction string `json:"idTransaction"`
	Data          []entity.Person `json:"data"`
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