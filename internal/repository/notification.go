package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type NotificationRepository interface {
	FindAllNotification(userID uuid.UUID) ([]entity.Notification, error)
	FindNotificationByID(notificationID uuid.UUID) (*entity.Notification, error)
	FindNotificationSeen(userID uuid.UUID) ([]entity.Notification, error)
	CreateNotification(notification *entity.Notification) (*entity.Notification, error)
	MarkAllNotificationsAsSeen(userID uuid.UUID) ([]entity.Notification, error)
	DeleteNotification(nnotification *entity.Notification) (bool, error)
	DeleteSeenAllNotification(userID uuid.UUID) (bool, error)
}
type notificationRepository struct {
	db *gorm.DB
	cacheable cache.Cacheable
}

func NewNotificationRepository(db *gorm.DB, cacheable cache.Cacheable) NotificationRepository {
	return &notificationRepository{db: db, cacheable: cacheable}
}


func (r *notificationRepository) FindAllNotification(userID uuid.UUID) ([]entity.Notification, error) {

	notifications := make([]entity.Notification, 0)

	if err := r.db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return notifications, err
	}
	return notifications, nil
}

func (r *notificationRepository) FindNotificationByID(notificationID uuid.UUID) (*entity.Notification, error) {

	notification := new(entity.Notification)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Where("id = ?", notificationID).Take(notification).Error; err != nil {
			return  err
		}
		
		if notification.IsSeen == false {
			if err := r.db.Model(notification).Where("id = ?", notificationID).Update("is_seen", true).Error; err != nil {
				return err
			}
		}

		return nil
	})


	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (r *notificationRepository) CreateNotification(notification *entity.Notification) (*entity.Notification, error) {

	if err := r.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func (r *notificationRepository) MarkAllNotificationsAsSeen(userID uuid.UUID) ([]entity.Notification, error) {
    

	if err := r.db.Model(&entity.Notification{}).Where("user_id = ?", userID).Update("is_seen", true).Error; err != nil {
		return nil, err
	}

	 data,err := r.FindAllNotification(userID)

	 if err != nil {
		return nil, err
	 }

	 return data, nil
}

func (r *notificationRepository) DeleteSeenAllNotification(userID uuid.UUID) (bool, error) {

	if err := r.db.Where("user_id = ? AND is_seen = ?", userID, true).Delete(&entity.Notification{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *notificationRepository) FindNotificationSeen(userID uuid.UUID) ([]entity.Notification, error) {

	notifications := make([]entity.Notification, 0)

	if err := r.db.Where("user_id = ? AND is_seen = ?", userID, true).Find(&notifications).Error; err != nil {
		return notifications, err
	}

	return notifications, nil
}

func (r *notificationRepository) DeleteNotification(notification *entity.Notification) (bool, error) {

	if err := r.db.Delete(&notification).Error; err != nil {
		return false, err
	}

	return true, nil
}