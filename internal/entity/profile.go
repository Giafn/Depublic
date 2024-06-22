package entity

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FullName       string    `json:"full_name"`
	Gender         string    `json:"gender"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	PhoneNumber    string    `json:"phone_number"`
	ProfilePicture string    `json:"profile_picture"`
	UserID         uuid.UUID `json:"user_id"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	Auditable
}

func NewProfile(fullName string, gender string, dateOfBirth time.Time, phoneNumber string, userID uuid.UUID, city string, province string) *Profile {
	return &Profile{
		ID:          uuid.New(),
		FullName:    fullName,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
		PhoneNumber: phoneNumber,
		UserID:      userID,
		City:        city,
		Province:    province,
		Auditable:   NewAuditable(),
	}
}

func UpdateProfile(userID uuid.UUID, fullName string, gender string, dateOfBirth time.Time, phoneNumber string, city string, province string) *Profile {
	return &Profile{
		UserID:      userID,
		FullName:    fullName,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
		PhoneNumber: phoneNumber,
		Auditable:   UpdateAuditable(),
	}
}
