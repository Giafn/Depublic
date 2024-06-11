package service

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(email string, password string) (jwtResponse, error)
	CreateUser(user *entity.User) (*entity.User, error)
	FindAllUser() ([]entity.User, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	tokenUseCase   token.TokenUseCase
}

type jwtResponse struct {
	Token      string `json:"token"`
	Expired_at string `json:"expired_at"`
}

func NewUserService(userRepository repository.UserRepository, tokenUseCase token.TokenUseCase) UserService {
	return &userService{
		userRepository: userRepository,
		tokenUseCase:   tokenUseCase,
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

	claims := token.JwtCustomClaims{
		ID:    user.UserId.String(),
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "Depublic-App",
		},
	}

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
