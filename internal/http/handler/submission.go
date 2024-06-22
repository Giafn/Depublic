package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type submissionHandler struct {
	submissionService service.SubmissionService
}

type SubmissionHandler interface {
	CreateSubmission(c echo.Context) error
	ListSubmission(c echo.Context) error
	AcceptSubmission(c echo.Context) error
	RejectSubmission(c echo.Context) error
	GetSubmissionByID(c echo.Context) error
}

func NewSubmissionHandler(submissionService service.SubmissionService) *submissionHandler {
	return &submissionHandler{submissionService}
}

func (h *submissionHandler) CreateSubmission(c echo.Context) error {
	input := binder.CreateSubmission{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	TransactionID, err := uuid.Parse(input.TransactionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	transaction, err := h.submissionService.FindTransactionByID(TransactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to find transaction"))
	}
	user, err := h.submissionService.FindUserByID(transaction.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to find transaction user"))
	}
	event, err := h.submissionService.FindEventByID(transaction.EventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to find event"))
	}

	if !event.MustUploadSubmission {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Event tidak memerlukan submission"))
	}

	_, err = h.submissionService.FindSubmissionByTransactionID(TransactionID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Submission sudah dibuat sebelumnya"))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "File tidak ditemukan"))
	}
	if file == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "File tidak ditemukan"))
	}

	filename, err := h.submissionService.UploadSubmission(file)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	_, err = h.submissionService.CreateSubmission(&entity.Submission{
		EventID:       transaction.EventID,
		UserID:        transaction.UserID,
		TransactionID: TransactionID,
		Name:          user.Profiles.FullName,
		Status:        "pending",
		Type:          "-",
		Filename:      filename,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse(http.StatusCreated, "Submission created", nil))
}

func (h *submissionHandler) ListSubmission(c echo.Context) error {
	return nil
}

func (h *submissionHandler) AcceptSubmission(c echo.Context) error {
	return nil
}

func (h *submissionHandler) RejectSubmission(c echo.Context) error {
	return nil
}

func (h *submissionHandler) GetSubmissionByID(c echo.Context) error {
	return nil
}
