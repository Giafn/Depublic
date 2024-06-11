package entity

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	Auditable
}

func NewUser(email, password, role string) *User {
	return &User{
		UserId:    uuid.New(),
		Email:     email,
		Password:  password,
		Role:      role,
		Auditable: NewAuditable(),
	}
}

func UpdateUser(user_id uuid.UUID, email, password, role, alamat, noHp string) *User {
	return &User{
		UserId:    user_id,
		Email:     email,
		Password:  password,
		Role:      role,
		Auditable: UpdateAuditable(),
	}
}
