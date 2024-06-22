package binder

type NotificationByID struct {
	ID string `param:"id" validate:"required"`
}

type CreateNotification struct {
	UserID  string `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}