package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	Auditable
}

func NewUser(email, password, role string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		Role:      role,
		Auditable: NewAuditable(),
	}
}

func UpdateUser(id uuid.UUID, email, password, role, alamat, noHp string) *User {
	return &User{
		ID:        id,
		Email:     email,
		Password:  password,
		Role:      role,
		Auditable: UpdateAuditable(),
	}
}
