package handler

import (
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
	var input binder.EventCreateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	 if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	} 

	
	startTime, _ := time.Parse("2006-01-02 15:04:05", input.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", input.EndTime)
	event := entity.NewEvent(input.Name, input.Organizer, input.Description, startTime, endTime, input.MustUploadSubmission, input.Province, input.City, input.District, input.FullAddress, input.Latitude, input.Longitude)
	pricings := make([]entity.Pricing, 0)
	for _, p := range input.Pricings {
		data := entity.Pricing{
			PricingId:        uuid.New(),
			EventID:   event.ID,
			Name:      p.Name,
			Quota:     p.Quota,
			Remaining: p.Remaining,
			Fee:       p.Fee,
			Auditable: entity.NewAuditable(),
		}
		pricings = append(pricings, data)
	}
	


	respData, err := h.eventService.CreateEvent(event, pricings);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses membuat event", respData))
}
