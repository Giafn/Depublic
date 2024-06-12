package service

import (
	"mime/multipart"
	"strings"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/encrypt"
	"github.com/Giafn/Depublic/pkg/upload"
	"github.com/google/uuid"
)

type profileService struct {
	profileRepo repository.ProfileRepository
	encryptTool    encrypt.EncryptTool
}

type ProfileService interface {
	FindProfileByUserID(userID uuid.UUID) (*entity.Profile, error)
	UpdateProfile(profile *entity.Profile,  file *multipart.FileHeader) (*entity.Profile, error)
	CreateProfile(profile *entity.Profile,  file *multipart.FileHeader) (*entity.Profile, error)
	Deleteprofile(userID uuid.UUID) (bool, error)
}

func NewProfileService(profilRepo repository.ProfileRepository, encryptTool encrypt.EncryptTool) ProfileService {
	return &profileService{profileRepo: profilRepo, encryptTool: encryptTool}
}

func (s *profileService) FindProfileByUserID(userID uuid.UUID) (*entity.Profile, error) {
	profile, err := s.profileRepo.FindProfileByUserID(userID)

	if err != nil {
		return nil, err
	}

	profile.PhoneNumber, _ = s.encryptTool.Decrypt(profile.PhoneNumber)
	return profile, nil
}

func (s *profileService) UpdateProfile(profile *entity.Profile,  file *multipart.FileHeader) (*entity.Profile, error) {

	oldProfile, err := s.profileRepo.FindProfileByUserID(profile.UserID)

	if err != nil {
		return nil, err
	}

	if file != nil {
		if err := upload.DeleteFile(oldProfile.ProfilePicture); err != nil {
			return nil, err
		}

		profilePic, err := upload.UploadImage(file, "profile-img")

		if err != nil {
			return nil, err
		}

		profile.ProfilePicture = profilePic
	} 

	if profile.Gender != "" {
		profile.Gender = strings.ToUpper(profile.Gender)
	}
	
	if profile.PhoneNumber != "" {
		profile.PhoneNumber, _ = s.encryptTool.Encrypt(profile.PhoneNumber)
	}
	
	
	newProfile, err := s.profileRepo.UpdateProfile(profile)

	if err != nil {
		return nil, err
	}

	newProfile.PhoneNumber, _ = s.encryptTool.Decrypt(newProfile.PhoneNumber)
	
	return newProfile, nil
}

func (s *profileService) CreateProfile(profile *entity.Profile,  file *multipart.FileHeader) (*entity.Profile, error) {

	if file != nil {
		profilePic, err := upload.UploadImage(file, "profile-img")
		if err != nil {
			return nil, err
		}
		profile.ProfilePicture = profilePic
	}


	profile.Gender = strings.ToUpper(profile.Gender)
	profile.PhoneNumber, _ = s.encryptTool.Encrypt(profile.PhoneNumber)
	newProfile, err := s.profileRepo.CreateProfile(profile)

	if err != nil {
		return nil, err
	}

	newProfile.PhoneNumber, _ = s.encryptTool.Decrypt(newProfile.PhoneNumber)

	return newProfile, nil
	
}

func (s *profileService) Deleteprofile(userID uuid.UUID) (bool, error) {
	profile, err := s.profileRepo.FindProfileByUserID(userID)

	if err != nil {
		return false, err
	}

	if err := upload.DeleteFile(profile.ProfilePicture); err != nil {
		return false, err
	}

	return s.profileRepo.DeleteProfile(profile)
}