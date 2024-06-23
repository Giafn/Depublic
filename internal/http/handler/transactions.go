package handler

import (
	"encoding/json"
	"fmt"
	"io"
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

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	input := new(binder.TransactionCreateRequest)

	if err := c.Bind(input); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	if len(input.Tickets) != input.TicketQuantity {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ticket quantity"))
	}

	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	transaction, mustUpload, err := h.transactionService.CreateTransaction(input.EventID, userID, input.Tickets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	res := map[string]string{
		"total_amount":           fmt.Sprintf("%d", transaction.TotalAmount),
		"must_upload_submission": "false",
		"payment_url":            "",
		"transaction_id":         transaction.ID.String(),
	}

	if mustUpload {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Transaction created successfully, please upload your submission", map[string]string{
			"must_upload_submission": "true",
			"transaction_id":         transaction.ID.String(),
		}))
	}

	res["payment_url"] = transaction.PaymentURL

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Transaction created successfully", res))
}

func (h *TransactionHandler) FindTransactionByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	transaction, err := h.transactionService.FindTransactionByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Transaction not found"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data transaksi", transaction))
}

func (h *TransactionHandler) FindAllTransactions(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	transactions, count, err := h.transactionService.FindAllTransactions(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// pagination
	data := pkg.Paginate(transactions, count, page, limit)

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data transaksi", data))
}

func (h *TransactionHandler) FindMyTransactions(c echo.Context) error {
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

	transactions, count, err := h.transactionService.FindMyTransactions(userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	data := pkg.Paginate(transactions, count, page, limit)

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data transaksi", data))
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	var input binder.TransactionUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	transaction := &entity.Transaction{
		ID:             id,
		EventID:        input.EventID,
		UserID:         input.UserID,
		TicketQuantity: input.TicketQuantity,
		TotalAmount:    input.TotalAmount,
		IsPaid:         input.IsPaid,
	}

	updatedTransaction, err := h.transactionService.UpdateTransaction(transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses mengupdate transaksi", updatedTransaction))
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	err = h.transactionService.DeleteTransaction(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Transaction deleted successfully", nil))
}

func (h *TransactionHandler) WebhookPayment(c echo.Context) error {
	fmt.Println("Received webhook payload")

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	request := &binder.MidtransWebhookRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	transID, _ := uuid.Parse(request.OrderID)

	transaction, err := h.transactionService.FindTransactionByID(transID)
	if err != nil {
		fmt.Println("Error finding transaction:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	user, err := h.transactionService.GetUsersById(transaction.UserID)
	if err != nil {
		fmt.Println("Error finding user:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	transactionStatus := request.TransactionStatus
	if transactionStatus != "settlement" && transactionStatus != "capture" {
		fmt.Println("Transaction status not accepted")

		transaction.Status = transactionStatus
		if transactionStatus == "pending" {
			transaction.Status = "pending_or_expired"
		}
		_, err = h.transactionService.UpdateTransaction(transaction)
		if err != nil {
			fmt.Println("Error updating transaction:", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		if transactionStatus != "pending" {
			html := service.CreateSuccessPaymentEmailHtml(user.Profiles.FullName, transaction.Status, transaction.TotalAmount, transaction.ID.String())
			service.ScheduleEmails(
				user.Email,
				"Payment Failed",
				html,
			)

			notif := entity.Notification{
				UserID:  transaction.UserID,
				Title:   "Payment Failed",
				Content: fmt.Sprintf("Your payment for transaction %s has been failed", transaction.ID.String()),
			}

			_, err = h.transactionService.CreateNotification(&notif)
			if err != nil {
				fmt.Println("Error creating notification:", err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}

		return c.String(http.StatusOK, "Webhook received successfully")
	}

	if transaction.IsPaid {
		fmt.Println("Transaction already paid")
		return c.String(http.StatusOK, "Webhook received successfully")
	}

	tickets, err := h.transactionService.GetTicketsByTransactionID(transID)
	if err != nil {
		fmt.Println("Error getting tickets:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	err = h.transactionService.UpdateTicketRemaining(tickets)
	if err != nil {
		fmt.Println("Error checking ticket availability:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	transaction.IsPaid = true
	transaction.Status = "paid"
	_, err = h.transactionService.UpdateTransaction(transaction)

	if err != nil {
		fmt.Println("Error updating transaction:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	html := service.CreateSuccessPaymentEmailHtml(user.Profiles.FullName, "accepted", transaction.TotalAmount, transaction.ID.String())
	service.ScheduleEmails(
		user.Email,
		"Payment Confirmation",
		html,
	)

	notif := entity.Notification{
		UserID:  transaction.UserID,
		Title:   "Payment Confirmation",
		Content: fmt.Sprintf("Your payment for transaction %s has been confirmed", transaction.ID.String()),
	}

	_, err = h.transactionService.CreateNotification(&notif)
	if err != nil {
		fmt.Println("Error creating notification:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, "Webhook received successfully")
}

func (h *TransactionHandler) PaymentRedirect(c echo.Context) error {
	payId := c.QueryParam("pay_id")
	decryptedURL, err := h.transactionService.DecryptPaymentURL(payId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	transID := c.QueryParam("transaction_id")

	// check is if submission not uploaded
	transaction, err := h.transactionService.FindTransactionByID(uuid.MustParse(transID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	event, err := h.transactionService.GetEventByID(transaction.EventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	if event.MustUploadSubmission {
		submission, err := h.transactionService.GetSubmissionByTransactionID(uuid.MustParse(transID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
		if submission == nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Your Submission not uploaded"))
		}
		if submission.Status != "accepted" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Your Submission not approved"))
		}
	}

	isAvaliable, err := h.transactionService.CheckTicketAvailability(uuid.MustParse(transID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if !isAvaliable {
		transaction.Status = "ticket_not_available"
		_, err = h.transactionService.UpdateTransaction(transaction)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Ticket is not available"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, decryptedURL)
}
