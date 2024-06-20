package binder

import (
	"time"

	"github.com/Giafn/Depublic/internal/entity"
)

type EventCreateRequest struct {
	Name                 string           `json:"name" validate:"required"`
	Organizer            string           `json:"organizer" validate:"required"`
	Description          string           `json:"description" validate:"required"`
	StartTime            string           `json:"start_time" validate:"required"`
	EndTime              string           `json:"end_time" validate:"required"`
	MustUploadSubmission bool             `json:"must_upload_submission" validate:"required"`
	Province             string           `json:"province" validate:"required"`
	City                 string           `json:"city" validate:"required"`
	District             string           `json:"district" validate:"required"`
	FullAddress          string           `json:"full_address" validate:"required"`
	Latitude             float64          `json:"latitude" validate:"required"`
	Longitude            float64          `json:"longitude" validate:"required"`
	Pricings             []PricingRequest `json:"pricings"`
}

type PricingRequest struct {
	Name      string `json:"name" validate:"required"`
	Quota     int    `json:"quota" validate:"required"`
	Remaining int    `json:"remaining" validate:"required"`
	Fee       int    `json:"fee" validate:"required"`
}

type EventUpdateRequest struct {
	Name                 string           `json:"name"`
	Organizer            string           `json:"organizer"`
	StartTime            time.Time        `json:"start_time"`
	EndTime              time.Time        `json:"end_time"`
	MustUploadSubmission bool             `json:"must_upload_submission"`
	Province             string           `json:"province"`
	City                 string           `json:"city"`
	District             string           `json:"district"`
	FullAddress          string           `json:"full_address"`
	Latitude             float64          `json:"latitude"`
	Longitude            float64          `json:"longitude"`
	Pricings             []entity.Pricing `json:"pricings"`
}
