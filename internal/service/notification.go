package service

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type NotificationService interface {
	FindAllNotification(userID uuid.UUID) ([]entity.Notification, error)
	FindNotificationByID(notificationID uuid.UUID) (*entity.Notification, error)
	CreateNotification(notification *entity.Notification) (*entity.Notification, error)
	UpdateSeenAllNotification(userID uuid.UUID) ([]entity.Notification, error)
	DeleteNotification(notificationID uuid.UUID) (bool, error)
	DeleteSeenNotification(userID uuid.UUID) (bool, error)
}

type notificationService struct {
	NotificationRepo repository.NotificationRepository
}


func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{NotificationRepo: notificationRepo}
}


func (s *notificationService) FindAllNotification(userID uuid.UUID) ([]entity.Notification, error) {
	return s.NotificationRepo.FindAllNotification(userID)
}

func (s *notificationService) FindNotificationByID(notificationID uuid.UUID) (*entity.Notification, error) {
	return s.NotificationRepo.FindNotificationByID(notificationID)
}

func (s *notificationService) CreateNotification(notification *entity.Notification) (*entity.Notification, error) {
	return s.NotificationRepo.CreateNotification(notification)
}

func (s *notificationService) UpdateSeenAllNotification(userID uuid.UUID) ([]entity.Notification, error) {
	return s.NotificationRepo.UpdateSeenAllNotification(userID)
}

func (s *notificationService) DeleteNotification(notificationID uuid.UUID) (bool, error) {
    notification, err := s.NotificationRepo.FindNotificationByID(notificationID)

	if err != nil {
		return false, err
	}

	return s.NotificationRepo.DeleteNotification(notification)
}

func (s *notificationService) DeleteSeenNotification(userID uuid.UUID) (bool, error) {
	return s.NotificationRepo.DeleteSeenAllNotification(userID)
}

