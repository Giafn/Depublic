package binder

type EventCreateRequest struct {
	Name                 string           `json:"name" validate:"required"`
	Organizer            string           `json:"organizer" validate:"required"`
	Description          string           `json:"description" validate:"required"`
	StartTime            string           `json:"start_time" validate:"required,date_with_time"`
	EndTime              string           `json:"end_time" validate:"required,date_with_time"`
	MustUploadSubmission bool             `json:"must_upload_submission"`
	Province             string           `json:"province" validate:"required"`
	City                 string           `json:"city" validate:"required"`
	District             string           `json:"district" validate:"required"`
	FullAddress          string           `json:"full_address" validate:"required"`
	Latitude             float64          `json:"latitude" validate:"required"`
	Longitude            float64          `json:"longitude" validate:"required"`
	Pricings             []PricingRequest `json:"pricings" validate:"required,dive"`
}

type PricingCreateRequest struct {
	EventID   string `json:"event_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Quota     int    `json:"quota" validate:"required"`
	Remaining int    `json:"remaining" validate:"required"`
	Fee       int    `json:"fee" validate:"required"`
}

type PricingRequest struct {
	Name      string `json:"name" validate:"required"`
	Quota     int    `json:"quota" validate:"required"`
	Remaining int    `json:"remaining" validate:"required"`
	Fee       int    `json:"fee" validate:"required"`
}

type PricingUpdateRequest struct {
	PricingId string `json:"pricing_id"`
	Name      string `json:"name"`
	Quota     int    `json:"quota"`
	Remaining int    `json:"remaining"`
	Fee       int    `json:"fee"`
}

type EventFindById struct {
	ID string `param:"id" validate:"required"`
}
type PricingFindById struct {
	ID string `param:"id" validate:"required"`
}

type DistanceRequest struct {
	Latitude float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type EventUpdateRequest struct {
	ID                   string                 `param:"id" validate:"required"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Organizer            string                 `json:"organizer"`
	StartTime            string                 `json:"start_time"`
	EndTime              string                 `json:"end_time"`
	MustUploadSubmission bool                   `json:"must_upload_submission"`
	Province             string                 `json:"province"`
	City                 string                 `json:"city"`
	District             string                 `json:"district"`
	FullAddress          string                 `json:"full_address"`
	Latitude             float64                `json:"latitude"`
	Longitude            float64                `json:"longitude"`
	Pricings             []PricingUpdateRequest `json:"pricings"`
}
