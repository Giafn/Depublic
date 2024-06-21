package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	NotificationService service.NotificationService
}


func NewNotificationHandler(notificationService service.NotificationService) NotificationHandler {
	return NotificationHandler{
		NotificationService: notificationService,
	}
}


func (h *NotificationHandler) FindAllNotification(c echo.Context) error {
	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	notifications, err := h.NotificationService.FindAllNotification(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan seluruh notifikasi", notifications))
}


func (h *NotificationHandler) FindNotificationByID(c echo.Context) error {

    var input binder.NotificationByID

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}
	
	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	notificationID := uuid.MustParse(input.ID)

	notification, err := h.NotificationService.FindNotificationByID(notificationID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mengecek notifikasi", notification))
}

func (h *NotificationHandler) DeleteNotificationByID(c echo.Context) error {

	var input binder.NotificationByID

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}
	
	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	notificationID := uuid.MustParse(input.ID)

	isDeleted, err := h.NotificationService.DeleteNotification(notificationID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menghapus notifikasi", isDeleted))
}

func (h *NotificationHandler) DeleteSeenNotifications(c echo.Context) error {
	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	isDeleted, err := h.NotificationService.DeleteSeenNotification(userID)

	if err != nil {	
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menghapus notifikasi yang sudah dilihat", isDeleted))
}

func (h *NotificationHandler) UpdateSeenNotifications(c echo.Context) error {
	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	isUpdated, err := h.NotificationService.UpdateSeenAllNotification(userID)

	if err != nil {	
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mengupdate semua notifikasi menjadi sudah dilihat", isUpdated))
}
