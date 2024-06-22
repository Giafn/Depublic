package entity

import (
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsSeen    bool      `json:"is_seen"`
	Auditable
}


func NewNotification(userID uuid.UUID, title string, content string) *Notification {
	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		IsSeen:    false,
		Auditable: NewAuditable(),
	}
}
