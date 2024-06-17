package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type profileRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

type ProfileRepository interface {
	CreateProfile(profile *entity.Profile) (*entity.Profile, error)
	UpdateProfile(profile *entity.Profile) (*entity.Profile, error)
	FindProfileByUserID(userID uuid.UUID) (*entity.Profile, error)
	DeleteProfile(profile *entity.Profile) (bool, error)
}

func NewProfileRepository(db *gorm.DB, cacheable cache.Cacheable) ProfileRepository {
	return &profileRepository{db: db, cacheable: cacheable}
}

func (r *profileRepository) CreateProfile(profile *entity.Profile) (*entity.Profile, error) {

	if err := r.db.Create(&profile).Error; err != nil {
		return profile, err
	}

	return profile, nil
}

func (r *profileRepository) UpdateProfile(profile *entity.Profile) (*entity.Profile, error) {

	key := fmt.Sprintf("Profile-%s", profile.UserID.String())

	fields := make(map[string]interface{})

	if profile.FullName != "" {
		fields["full_name"] = profile.FullName
	}

	if profile.Gender != "" {
		fields["gender"] = profile.Gender
	}

	if profile.DateOfBirth != (time.Time{}) {
		fields["date_of_birth"] = profile.DateOfBirth
	}

	if profile.PhoneNumber != "" {
		fields["phone_number"] = profile.PhoneNumber
	}

	if profile.ProfilePicture != "" {
		fields["profile_picture"] = profile.ProfilePicture
	}

	if err := r.db.Model(&profile).Where("user_id = ?", profile.UserID).Updates(fields).Error; err != nil {
		return profile, err
	}

	if data := r.cacheable.Get(key); data != "" {
		err := r.cacheable.Del(key)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}

func (r *profileRepository) FindProfileByUserID(userID uuid.UUID) (*entity.Profile, error) {

	key := fmt.Sprintf("Profile-%s", userID.String())

	profile := new(entity.Profile)

	data := r.cacheable.Get(key)

	if data == "" {
		if err := r.db.Where("user_id = ?", userID).Take(&profile).Error; err != nil {
			return nil, err
		}
		marshalledProfile, _ := json.Marshal(profile)
		err := r.cacheable.Set(key, marshalledProfile, 2*time.Minute)
		if err != nil {
			return nil, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &profile)
		if err != nil {
			return profile, err
		}
	}

	return profile, nil
}

func (r *profileRepository) DeleteProfile(profile *entity.Profile) (bool, error) {

	key := fmt.Sprintf("Profile-%s", profile.UserID.String())

	if err := r.db.Delete(&profile).Error; err != nil {
		return false, err
	}

	if data := r.cacheable.Get(key); data != "" {
		err := r.cacheable.Del(key)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
