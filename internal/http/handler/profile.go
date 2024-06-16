package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type profileHandler struct {
	profileService service.ProfileService
}

type ProfileHandler interface {
	FindCurrentUserProfile(c echo.Context) error
	UpdateProfile(c echo.Context) error
	DeleteProfile(c echo.Context) error
}

func NewProfileHandler(profileService service.ProfileService) ProfileHandler {
	return &profileHandler{
		profileService: profileService,
	}
}

func (h *profileHandler) FindCurrentUserProfile(c echo.Context) error {

	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	profile, err := h.profileService.FindProfileByUserID(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data profile", map[string]interface{}{
		"id":              profile.ID,
		"full_name":       profile.FullName,
		"gender":          profile.Gender,
		"date_of_birth":   profile.DateOfBirth,
		"phone_number":    profile.PhoneNumber,
		"profile_picture": profile.ProfilePicture,
		"city":            profile.City,
		"province":        profile.Province,
		"created_at":      profile.CreatedAt,
	}))
}

func (h *profileHandler) UpdateProfile(c echo.Context) error {
	var input binder.UpdateProfileRequest

	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	file, err := c.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	dateOfBirth, _ := time.Parse("2006-01-02", input.DateOfBirth)
	fmt.Println(dateOfBirth)
	updatedProfile := entity.UpdateProfile(userID, input.FullName, input.Gender, dateOfBirth, input.PhoneNumber, input.City, input.Province)

	profile, err := h.profileService.UpdateProfile(updatedProfile, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	type UpdateProfileResponse struct {
		FullName       string    `json:"full_name,omitempty"`
		Gender         string    `json:"gender,omitempty"`
		DateOfBirth    time.Time `json:"date_of_birth,omitempty"`
		PhoneNumber    string    `json:"phone_number,omitempty"`
		ProfilePicture string    `json:"profile_picture,omitempty"`
		City           string    `json:"city,omitempty"`
		Province       string    `json:"province,omitempty"`
		UpdateAt       time.Time `json:"update_at"`
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses update profile", UpdateProfileResponse{
		FullName:       profile.FullName,
		Gender:         profile.Gender,
		DateOfBirth:    profile.DateOfBirth,
		PhoneNumber:    profile.PhoneNumber,
		ProfilePicture: profile.ProfilePicture,
		City:           profile.City,
		Province:       profile.Province,
		UpdateAt:       profile.UpdatedAt,
	}))
}

func (h *profileHandler) DeleteProfile(c echo.Context) error {

	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	isDeleted, err := h.profileService.Deleteprofile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete profile", isDeleted))
}
