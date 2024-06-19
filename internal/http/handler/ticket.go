package handler

import (
	"errors"
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"fmt"
)

type TicketHandler struct {
	ticketService service.TicketService
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

	idTransaction, err := uuid.Parse(input.IDTransaction)
	fmt.Println(idTransaction)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	_, err = h.transactionService.FindTransactionByID(idTransaction)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	for i := 0; i < len(input.Data); i++ {

		fmt.Println(input.Data[i].Name)
		newTicket := entity.NewTicket(input.IDTransaction, input.IDEvent, input.Data[i].Name)

		fmt.Println("id trx")
		fmt.Println(newTicket.IDTransaction)

		_, err := h.ticketService.CreateTicket(newTicket)

		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
		}

	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses membuat tiket", nil))
}

func (h *TicketHandler) Test(c echo.Context) error {
	fmt.Println("halo")

	return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "testtt aja emg sengaja error"))
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

	requestedTicket := entity.ValidateTicket(*oldTicket)

	validatedTicket, err := h.ticketService.ValidateTicket(requestedTicket)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, repository.ErrTicketAlreadyValidated) {
			return c.JSON(http.StatusConflict, response.ErrorResponse(http.StatusConflict, err.Error()))
		}
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// fmt.Println("3")
	// fmt.Println(validatedTicket.IsUsed)

	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	// }

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