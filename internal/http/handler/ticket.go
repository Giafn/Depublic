package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"fmt"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) TicketHandler {
	return TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) CreateTicket(c echo.Context) error {
	input := binder.TicketCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	fmt.Println(input.Data)

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

	id := uuid.MustParse(input.ID)

	ticket, err := h.ticketService.FindTicketByID(id)

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