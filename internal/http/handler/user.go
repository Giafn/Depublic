package handler

import (
	"fmt"
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	roleAdmin = "Admin"
	roleUser  = "User"
)

type UserHandler struct {
	userService service.UserService
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
	input := binder.UserRegisterRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	newUser := entity.NewUser(input.Email, input.Password, roleUser, false)

	user, err := h.userService.RegisterUser(newUser)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	data := map[string]interface{}{
		"email": user.Email,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mendaftar sebagai user", data))
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	input := binder.UserCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	if input.Role != roleAdmin && input.Role != roleUser {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "role tidak valid hanya menerima 'Admin' atau 'User'"))
	}

	newUser := entity.NewUser(input.Email, input.Password, input.Role, true)

	user, err := h.userService.CreateUser(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mendaftar sebagai user", user))
}

func (h *UserHandler) FindAllUser(c echo.Context) error {
	users, err := h.userService.FindAllUser()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data user", users))
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

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data user", user))
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
	tokenString = tokenString[7:]
	err := h.userService.Logout(tokenString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses logout", nil))
}
