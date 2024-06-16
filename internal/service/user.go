package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Giafn/Depublic/configs"
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/encrypt"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/Giafn/Depublic/pkg/upload"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(email string, password string) (jwtResponse, error)
	RegisterUser(user *binder.UserRegisterRequest, file *multipart.FileHeader) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	FindAllUser() ([]entity.User, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	VerifyEmail(id uuid.UUID) error
	ResendEmailVerification(email string) error
	Logout(tokenString string) error
}

type userService struct {
	userRepository    repository.UserRepository
	profileRepository repository.ProfileRepository
	tokenUseCase      token.TokenUseCase
	encryptTool       encrypt.EncryptTool
	cfg               *configs.Config
}

type jwtResponse struct {
	Token      string `json:"token"`
	Expired_at string `json:"expired_at"`
}

func NewUserService(
	userRepository repository.UserRepository,
	profileRepository repository.ProfileRepository,
	tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool,
	cfg *configs.Config,
) UserService {
	return &userService{
		userRepository:    userRepository,
		profileRepository: profileRepository,
		tokenUseCase:      tokenUseCase,
		encryptTool:       encryptTool,
		cfg:               cfg,
	}
}

func (s *userService) Login(email string, password string) (data jwtResponse, error error) {
	user, err := s.userRepository.FindUserByEmail(email)

	data = jwtResponse{}
	data.Token = ""
	data.Expired_at = ""

	if err != nil {
		return data, errors.New("email/password yang anda masukkan salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return data, errors.New("email/password yang anda masukkan salah")
	}

	if !user.IsVerified {
		return data, errors.New("silahkan verifikasi akun anda terlebih dahulu")
	}

	claims := s.tokenUseCase.CreateClaims(user.UserId.String(), user.Email, user.Role)

	token, expiredAt, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return data, err
	}

	data.Token = token
	data.Expired_at = expiredAt.Format("2006-01-02 15:04:05")

	return data, nil
}

func (s *userService) CreateUser(user *entity.User) (*entity.User, error) {
	_, err := s.userRepository.FindUserByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) RegisterUser(input *binder.UserRegisterRequest, file *multipart.FileHeader) (*entity.User, error) {

	user := entity.NewUser(input.Email, input.Password, "User", false)
	_, err := s.userRepository.FindUserByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	phone, _ := s.encryptTool.Encrypt(input.PhoneNumber)
	dateOfBirth, _ := time.Parse("2006-01-02", input.DateOfBirth)

	newProfile := entity.NewProfile(
		input.FullName,
		strings.ToUpper(input.Gender),
		dateOfBirth,
		phone,
		newUser.UserId,
		input.City,
		input.Province,
	)

	if file != nil {
		profilePic, err := upload.UploadImage(file, "profile-img")
		if err != nil {
			return nil, err
		}
		newProfile.ProfilePicture = profilePic
	}

	_, err = s.profileRepository.CreateProfile(newProfile)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%s/app/api/v1/account/verify/%s", s.cfg.Host, s.cfg.Port, newUser.UserId.String())
	html := "<h1>Account Confirmation</h1><p>Click <a href='" + url + "'>here</a> to confirm your account</p>"

	ScheduleEmails(
		user.Email,
		"Account Confirmation of Registration ",
		html,
	)

	return newUser, nil
}

func (s *userService) FindAllUser() ([]entity.User, error) {
	users, err := s.userRepository.FindAllUser()
	if err != nil {
		return nil, err
	}

	formattedUser := make([]entity.User, 0)
	for _, v := range users {
		formattedUser = append(formattedUser, entity.User{
			UserId:    v.UserId,
			Email:     v.Email,
			Role:      v.Role,
			Auditable: v.Auditable,
		})
	}

	return formattedUser, nil
}

func (s *userService) FindUserByID(id uuid.UUID) (*entity.User, error) {
	return s.userRepository.FindUserByID(id)
}

func (s *userService) VerifyEmail(id uuid.UUID) error {
	user, err := s.userRepository.FindUserByID(id)
	if err != nil {
		return err
	}

	user.IsVerified = true
	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ResendEmailVerification(email string) error {
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("akun anda sudah terverifikasi")
	}

	url := fmt.Sprintf("http://%s:%s/app/api/v1/account/verify/%s", s.cfg.Host, s.cfg.Port, user.UserId.String())
	html := "<h1>Account Confirmation</h1><p>Click <a href='" + url + "'>here</a> to confirm your account</p>"

	ScheduleEmails(
		user.Email,
		"Account Confirmation of Registration ",
		html,
	)

	return nil
}

func (s *userService) Logout(tokenString string) error {
	return s.tokenUseCase.InvalidateToken(tokenString)
}
