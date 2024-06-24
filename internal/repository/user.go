package repository

import (
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	CreateUserWithProfile(user *entity.User, profile *entity.Profile) (*entity.User, error)
	UpdateUserWithProfile(user *entity.User, profile *entity.Profile) (*entity.User, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindAllUser(page, limit int) ([]entity.User, int, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id uuid.UUID) error
}

type userRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewUserRepository(db *gorm.DB, cacheable cache.Cacheable) UserRepository {
	return &userRepository{db: db, cacheable: cacheable}
}

func (r *userRepository) FindUserByID(id uuid.UUID) (*entity.User, error) {
	user := new(entity.User)

	if err := r.db.Preload("Profiles").
		Where("users.user_id = ?", id).
		Take(user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	user := new(entity.User)
	if err := r.db.Preload("Profiles").
		Where("email = ?", email).
		Take(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindAllUser(page, limit int) ([]entity.User, int, error) {
	offset := (page - 1) * limit

	users := make([]entity.User, 0)
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Preload("Profiles").
		Find(&users).
		Error; err != nil {
		return users, 0, err
	}

	var count int64
	if err := r.db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return users, 0, err
	}

	return users, int(count), nil
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return user, err
	}
	err := r.cacheable.Del("FindAllUsers")
	if err != nil {
		return user, err
	}
	return user, nil
}

// create user with profile
func (r *userRepository) CreateUserWithProfile(user *entity.User, profile *entity.Profile) (*entity.User, error) {
	tx := r.db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	profile.UserID = user.UserId
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	tx.Commit()

	err := r.cacheable.Del("FindAllUsers")
	if err != nil {
		return user, err
	}

	return user, nil
}

// update user with profile
func (r *userRepository) UpdateUserWithProfile(user *entity.User, profile *entity.Profile) (*entity.User, error) {
	tx := r.db.Begin()

	if err := tx.Model(&entity.User{}).Where("user_id = ?", user.UserId).Updates(user).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	if err := tx.Model(&entity.Profile{}).Where("user_id = ?", user.UserId).Updates(profile).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	tx.Commit()

	err := r.cacheable.Del("FindAllUsers")
	if err != nil {
		return user, err
	}

	user.Profiles = *profile

	return user, nil
}

func (r *userRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	// Update the user in the database
	if err := r.db.Model(&entity.User{}).Where("user_id = ?", user.UserId).Updates(user).Error; err != nil {
		return user, err
	}

	err := r.cacheable.Del("FindAllUsers")
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	if err := r.db.Where("user_id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}

	// profile delete
	if err := r.db.Where("user_id = ?", id).Delete(&entity.Profile{}).Error; err != nil {
		return err
	}

	err := r.cacheable.Del("FindAllUsers")
	if err != nil {
		return err
	}

	return nil
}
