package service

import (
	"errors"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/encrypt"
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
	encryptTool    encrypt.EncryptTool
}

type jwtResponse struct {
	Token      string `json:"token"`
	Expired_at string `json:"expired_at"`
}

func NewUserService(userRepository repository.UserRepository, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool) UserService {
	return &userService{
		userRepository: userRepository,
		tokenUseCase:   tokenUseCase,
		encryptTool:    encryptTool,
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

	user.Alamat, _ = s.encryptTool.Decrypt(user.Alamat)
	user.NoHp, _ = s.encryptTool.Decrypt(user.NoHp)

	claims := token.JwtCustomClaims{
		ID:     user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		Alamat: user.Alamat,
		NoHP:   user.NoHp,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "Go-Commerce",
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	user.Alamat, err = s.encryptTool.Encrypt(user.Alamat)
	if err != nil {
		panic(err)
	}
	user.NoHp, err = s.encryptTool.Encrypt(user.NoHp)
	if err != nil {
		panic(err)
	}

	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	newUser.Alamat, _ = s.encryptTool.Decrypt(newUser.Alamat)
	newUser.NoHp, _ = s.encryptTool.Decrypt(newUser.NoHp)

	return newUser, nil
}

func (s *userService) FindAllUser() ([]entity.User, error) {
	users, err := s.userRepository.FindAllUser()
	if err != nil {
		return nil, err
	}

	formattedUser := make([]entity.User, 0)
	for _, v := range users {
		v.Alamat, _ = s.encryptTool.Decrypt(v.Alamat)
		v.NoHp, _ = s.encryptTool.Decrypt(v.NoHp)
		formattedUser = append(formattedUser, v)
	}

	return formattedUser, nil
}

func (s *userService) FindUserByID(id uuid.UUID) (*entity.User, error) {
	return s.userRepository.FindUserByID(id)
}
