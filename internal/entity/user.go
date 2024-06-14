package entity

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
	Auditable
}

func NewUser(email, password, role string, isActive bool) *User {
	return &User{
		UserId:    uuid.New(),
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  isActive,
		Auditable: NewAuditable(),
	}
}

func UpdateUser(user_id uuid.UUID, email, password, role string, isActive bool) *User {
	return &User{
		UserId:    user_id,
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  isActive,
		Auditable: UpdateAuditable(),
	}
}
