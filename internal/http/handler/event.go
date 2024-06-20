package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) EventHandler {
	return EventHandler{eventService}
}

func (h *EventHandler) CreateNewEvent(c echo.Context) error {
	input := new(binder.EventCreateRequest)
	fmt.Print(input)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	/* if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	} */

	pricings := make([]entity.Pricing, len(input.Pricings))
	for i, p := range input.Pricings {
		pricings[i] = entity.Pricing{
			ID:        uuid.New(),
			Name:      p.Name,
			Quota:     p.Quota,
			Remaining: p.Remaining,
			Fee:       p.Fee,
		}
	}
	startTime, _ := time.Parse("2006-01-02", input.StartTime)
	endTime, _ := time.Parse("2006-01-02", input.EndTime)
	event := entity.NewEvent(input.Name, input.Organizer, input.Description, startTime, endTime, input.MustUploadSubmission, input.Province, input.City, input.District, input.FullAddress, input.Latitude, input.Longitude)
	event.Pricings = pricings

	if err := h.eventService.CreateEvent(event); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, event)
}
