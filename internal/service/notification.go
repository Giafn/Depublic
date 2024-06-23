package service

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/google/uuid"
)

type NotificationService interface {
	FindAllNotification(userID uuid.UUID) ([]entity.Notification, error)
	FindNotificationByID(notificationID uuid.UUID) (*entity.Notification, error)
	CreateNotification(notification *entity.Notification) (*entity.Notification, error)
	MarkAllNotificationsAsSeen(userID uuid.UUID) ([]entity.Notification, error)
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

func (s *notificationService) MarkAllNotificationsAsSeen(userID uuid.UUID) ([]entity.Notification, error) {
	data, err := s.NotificationRepo.FindNotificationSeen(userID)

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("Tidak ada notifikasi yang sudah dilihat")
	}

	return s.NotificationRepo.MarkAllNotificationsAsSeen(userID)
}

func (s *notificationService) DeleteNotification(notificationID uuid.UUID) (bool, error) {
	notification, err := s.NotificationRepo.FindNotificationByID(notificationID)

	if err != nil || notification == nil {
		return false, err
	}

	return s.NotificationRepo.DeleteNotification(notification)
}

func (s *notificationService) DeleteSeenNotification(userID uuid.UUID) (bool, error) {
	data, err := s.NotificationRepo.FindNotificationSeen(userID)

	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return false, errors.New("Tidak ada notifikasi yang sudah dilihat")
	}

	return s.NotificationRepo.DeleteSeenAllNotification(userID)
}
