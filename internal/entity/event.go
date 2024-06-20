package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Organizer            string    `json:"organizer"`
	Description          string    `json:"description"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	MustUploadSubmission bool      `json:"must_upload_submission"`
	Province             string    `json:"province"`
	City                 string    `json:"city"`
	District             string    `json:"district"`
	FullAddress          string    `json:"full_address"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	Pricings             []Pricing `gorm:"foreignKey:EventID;references:ID" json:"pricings"`
	Auditable
}

func NewEvent(name, organizer, description string, startTime, endTime time.Time, mustUploadSubmission bool, province, city, district, fullAddress string, latitude, longitude float64) *Event {
	return &Event{
		ID:                   uuid.New(),
		Name:                 name,
		Organizer:            organizer,
		Description:          description,
		StartTime:            startTime,
		EndTime:              endTime,
		MustUploadSubmission: mustUploadSubmission,
		Province:             province,
		City:                 city,
		District:             district,
		FullAddress:          fullAddress,
		Latitude:             latitude,
		Longitude:            longitude,
		Auditable:            NewAuditable(),
	}
}

func UpdateEvent(id uuid.UUID, name, organizer string, startTime, endTime time.Time, mustUploadSubmission bool, province, city, district, fullAddress string, latitude, longitude float64, pricings []Pricing) *Event {
	return &Event{
		ID:                   id,
		Name:                 name,
		Organizer:            organizer,
		StartTime:            startTime,
		EndTime:              endTime,
		MustUploadSubmission: mustUploadSubmission,
		Province:             province,
		City:                 city,
		District:             district,
		FullAddress:          fullAddress,
		Latitude:             latitude,
		Longitude:            longitude,
		Pricings:             pricings,
		Auditable:            UpdateAuditable(),
	}
}
