package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	pkg "github.com/Giafn/Depublic/pkg/pagination"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
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
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)
	submissions, count, err := h.submissionService.ListSubmission(userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if len(submissions) == 0 {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Submission is empty", nil))
	}

	data := make([]binder.ListSubmissionResponse, 0)
	for _, submission := range submissions {
		data = append(data, binder.ListSubmissionResponse{
			ID:            submission.ID,
			EventID:       submission.EventID,
			UserID:        submission.UserID,
			TransactionID: submission.TransactionID,
			Name:          submission.Name,
			Status:        submission.Status,
			Type:          submission.Type,
			Filename:      submission.Filename,
		})
	}

	dataRes := pkg.Paginate(data, count, page, limit)

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "List submission", dataRes))

}

func (h *submissionHandler) AcceptSubmission(c echo.Context) error {
	submissionID := c.Param("id")
	id, err := uuid.Parse(submissionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	submission, err := h.submissionService.FindSubmissionByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Submission not found"))
	}

	if submission.Status == "accepted" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Submission already accepted"))
	}

	if submission.Status == "rejected" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Submission already rejected"))
	}

	submission.Status = "accepted"
	_, err = h.submissionService.UpdateSubmission(submission)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.submissionService.SendEmailSubmission(submission.Status, submission)
	if err != nil {
		fmt.Println("Failed to send email")
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Submission accepted", nil))
}

func (h *submissionHandler) RejectSubmission(c echo.Context) error {
	submissionID := c.Param("id")
	id, err := uuid.Parse(submissionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	submission, err := h.submissionService.FindSubmissionByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Submission not found"))
	}

	if submission.Status == "accepted" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Submission already accepted"))
	}

	submission.Status = "rejected"
	_, err = h.submissionService.UpdateSubmission(submission)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.submissionService.SendEmailSubmission(submission.Status, submission)
	if err != nil {
		fmt.Println("Failed to send email")
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Submission rejected", nil))
}

func (h *submissionHandler) GetSubmissionByID(c echo.Context) error {
	return nil
}
