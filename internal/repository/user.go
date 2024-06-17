package repository

import (
	"encoding/json"
	"time"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	CreateUserWithProfile(user *entity.User, profile *entity.Profile) (*entity.User, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindAllUser() ([]entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
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

	if err := r.db.Where("users.user_id = ?", id).
		Take(user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	user := new(entity.User)
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindAllUser() ([]entity.User, error) {
	users := make([]entity.User, 0)

	key := "FindAllUsers"

	data := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Find(&users).Error; err != nil {
			return users, err
		}
		marshalledUsers, _ := json.Marshal(users)
		err := r.cacheable.Set(key, marshalledUsers, 5*time.Minute)
		if err != nil {
			return users, err
		}
	} else {
		// Data ditemukan di Redis, unmarshal data ke users
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			return users, err
		}
	}

	return users, nil
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
