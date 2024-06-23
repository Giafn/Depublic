package handler

import (
	"errors"
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"fmt"
)

type TicketHandler struct {
	ticketService      service.TicketService
	transactionService service.TransactionService
}

func NewTicketHandler(ticketService service.TicketService, transactionService service.TransactionService) TicketHandler {
	return TicketHandler{ticketService: ticketService, transactionService: transactionService}
}

func (h *TicketHandler) CreateTicket(c echo.Context) error {
	input := binder.TicketCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	transactionId, err := uuid.Parse(input.TransactionID)
	fmt.Println(transactionId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	_, err = h.transactionService.FindTransactionByID(transactionId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	for i := 0; i < len(input.Data); i++ {

		fmt.Println(input.Data[i].Name)
		newTicket := entity.NewTicket(input.TransactionID, input.EventID, input.Data[i].Name, "")

		fmt.Println("id trx")
		fmt.Println(newTicket.TransactionID)

		_, err := h.ticketService.CreateTicket(newTicket)

		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
		}

	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses membuat tiket", nil))
}

func (h *TicketHandler) FindTicketByID(c echo.Context) error {
	input := binder.TicketFindByIdRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id, err := uuid.Parse(input.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	ticket, err := h.ticketService.FindTicketByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data tiket", ticket))
}

func (h *TicketHandler) FindTicketByBookingNumber(c echo.Context) error {
	input := binder.TicketFindByBookingNumberRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	ticket, err := h.ticketService.FindTicketByBookingNumber(input.BookingNumber)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan data tiket", ticket))
}

func (h *TicketHandler) UpdateTicket(c echo.Context) error {
	input := binder.TicketUpdateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	oldTicket, err := h.ticketService.FindTicketByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	requestedTicket := entity.UpdateTicket(*oldTicket, input.Name)

	updatedTicket, err := h.ticketService.UpdateTicket(requestedTicket)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses update tiket", updatedTicket))
}

func (h *TicketHandler) ValidateTicket(c echo.Context) error {
	input := binder.TicketValidateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	oldTicket, err := h.ticketService.FindTicketByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	fmt.Println("old ticket")
	fmt.Println(oldTicket.IsUsed) //false

	// check if ticket transaction is_paid true
	transaction, err := h.transactionService.FindTransactionByID(uuid.MustParse(oldTicket.TransactionID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if !transaction.IsPaid {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "transaksi belum dibayar"))
	}

	validatedTicket, err := h.ticketService.ValidateTicket(oldTicket)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, repository.ErrTicketAlreadyValidated) {
			return c.JSON(http.StatusConflict, response.ErrorResponse(http.StatusConflict, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses memvalidasi tiket", validatedTicket))
}

func (h *TicketHandler) DeleteTicketById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err = h.ticketService.DeleteTicketById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "tiket berhasil dihapus", nil))
}

func (h *TicketHandler) DeleteTicketByBookingNumber(c echo.Context) error {
	bookingNumber := c.Param("bookingNum")

	err := h.ticketService.DeleteTicketByBookingNumber(bookingNumber)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "tiket berhasil dihapus", nil))
}

func (h *TicketHandler) FindAllTickets(c echo.Context) error {
	tickets, err := h.ticketService.FindAllTickets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menampilkan seluruh tiket", tickets))
}

func (h *TicketHandler) FindTicketsByTransactionId(c echo.Context) error {
	param := c.Param("transactionId")
	transactionId, err := uuid.Parse(param)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID := uuid.MustParse(claims.ID)

	tickets, err := h.ticketService.FindTicketsByTransactionId(transactionId, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	message := fmt.Sprintf("sukses menampilkan seluruh tiket dengan transaction id %s", param)

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, message, tickets))
}

func (h *TicketHandler) FindTicketsByUser(c echo.Context) error {
	dataUser, _ := c.Get("user").(*jwt.Token)
	claims := dataUser.Claims.(*token.JwtCustomClaims)

	userID, err := uuid.Parse(claims.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	tickets, err := h.ticketService.FindTicketsByUser(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	message := fmt.Sprintf("sukses menampilkan seluruh tiket milik user %s", userID)

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, message, tickets))
}
