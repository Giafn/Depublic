package entity

import "github.com/google/uuid"

type User struct {
	UserId     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	Auditable
}

func NewUser(email, password, role string, IsVerified bool) *User {
	return &User{
		UserId:     uuid.New(),
		Email:      email,
		Password:   password,
		Role:       role,
		IsVerified: IsVerified,
		Auditable:  NewAuditable(),
	}
}

func UpdateUser(user_id uuid.UUID, email, password, role string, IsVerified bool) *User {
	return &User{
		UserId:     user_id,
		Email:      email,
		Password:   password,
		Role:       role,
		IsVerified: IsVerified,
		Auditable:  UpdateAuditable(),
	}
}
