package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

type userResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	IsVerified     bool      `json:"is_verified"`
	FullName       string    `json:"full_name"`
	Gender         string    `json:"gender"`
	DateOfBirth    string    `json:"date_of_birth"`
	PhoneNumber    string    `json:"phone_number"`
	ProfilePicture string    `json:"profile_picture"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (h *UserHandler) Login(c echo.Context) error {
	input := new(binder.UserLoginRequest)

	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	user, err := h.userService.Login(input.Email, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "login success", user))
}

func (h *UserHandler) Register(c echo.Context) error {
	var input binder.UserRegisterRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	file, err := c.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	user, err := h.userService.RegisterUser(&input, file)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	data := map[string]interface{}{
		"email": user.Email,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mendaftar sebagai user", data))
}

func (h *UserHandler) FindAllUser(c echo.Context) error {
	users, err := h.userService.FindAllUser()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// map user to response
	usersResponse := make([]userResponse, 0)
	for _, user := range users {
		userMap := userResponse{
			UserID:         user.UserId,
			Email:          user.Email,
			Role:           user.Role,
			IsVerified:     user.IsVerified,
			FullName:       user.Profiles.FullName,
			Gender:         user.Profiles.Gender,
			DateOfBirth:    user.Profiles.DateOfBirth.Format("2006-01-02"),
			PhoneNumber:    user.Profiles.PhoneNumber,
			ProfilePicture: user.Profiles.ProfilePicture,
			City:           user.Profiles.City,
			Province:       user.Profiles.Province,
			CreatedAt:      user.CreatedAt.String(),
			UpdatedAt:      user.UpdatedAt.String(),
		}

		usersResponse = append(usersResponse, userMap)
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data user", usersResponse))
}

func (h *UserHandler) FindUserByID(c echo.Context) error {
	var input binder.UserFindByIDRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	user, err := h.userService.FindUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	formattedUser := userResponse{
		UserID:         user.UserId,
		Email:          user.Email,
		Role:           user.Role,
		IsVerified:     user.IsVerified,
		FullName:       user.Profiles.FullName,
		Gender:         user.Profiles.Gender,
		DateOfBirth:    user.Profiles.DateOfBirth.Format("2006-01-02"),
		PhoneNumber:    user.Profiles.PhoneNumber,
		ProfilePicture: user.Profiles.ProfilePicture,
		City:           user.Profiles.City,
		Province:       user.Profiles.Province,
		CreatedAt:      user.CreatedAt.String(),
		UpdatedAt:      user.UpdatedAt.String(),
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data user", formattedUser))
}

func (h *UserHandler) VerifyEmail(c echo.Context) error {
	var input binder.UserFindByIDRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	err := h.userService.VerifyEmail(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses verifikasi email", nil))
}

func (h *UserHandler) ResendEmailVerification(c echo.Context) error {
	input := new(binder.UserVerifyEmailRequest)

	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	email := input.Email

	err := h.userService.ResendEmailVerification(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mengirim ulang email verifikasi", nil))
}

func (h *UserHandler) Logout(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	} else {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid authorization header format"))
	}

	err := h.userService.Logout(tokenString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses logout", nil))
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var input binder.UserCreateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	file, err := c.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	user, err := h.userService.CreateUser(&input, file)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	data := map[string]interface{}{
		"user_id": user.UserId,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses membuat user", data))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	var input binder.UserUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	file, err := c.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	user, err := h.userService.UpdateUser(id, &input, file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userResponse := userResponse{
		UserID:         user.UserId,
		Email:          user.Email,
		Role:           user.Role,
		IsVerified:     user.IsVerified,
		FullName:       user.Profiles.FullName,
		Gender:         user.Profiles.Gender,
		DateOfBirth:    user.Profiles.DateOfBirth.Format("2006-01-02"),
		PhoneNumber:    user.Profiles.PhoneNumber,
		ProfilePicture: user.Profiles.ProfilePicture,
		City:           user.Profiles.City,
		Province:       user.Profiles.Province,
		CreatedAt:      user.CreatedAt.String(),
		UpdatedAt:      user.UpdatedAt.String(),
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses update user", userResponse))
}
