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
			PricingId: uuid.New(),
			EventID:   event.ID,
			Name:      p.Name,
			Quota:     p.Quota,
			Remaining: p.Remaining,
			Fee:       p.Fee,
			Auditable: entity.NewAuditable(),
		}
		pricings = append(pricings, data)
	}

	respData, err := h.eventService.CreateEvent(event, pricings)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses membuat event", respData))
}

func (h *EventHandler) FindEventByID(c echo.Context) error {
	var input binder.EventFindById

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input Invalid"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	event, err := h.eventService.FindEventByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Sukses menampilkan data event", event))
}

func (h *EventHandler) FindPricingByEventID(c echo.Context) error {
	var input binder.EventFindById

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input Invalid"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	id := uuid.MustParse(input.ID)

	price, err := h.eventService.FindPricingByEventID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Sukses menampilkan data event", price))
}

func (h *EventHandler) GetEvents(c echo.Context) error {
	allowedParams := map[string]bool{
		"price":    true,
		"province": true,
		"timeStart":     true,
		"category": true,
		"sort":     true,
	}

	allowedSorts := map[string]bool{
		"terbaru":  true,
		"terdekat": true,
		"termahal": true,
		"termurah": true,
	}

	allowedTime := map[string]bool{
		"day":true,
		"night": true,
	}

	allowedCategory := map[string]bool{
		"withSubmission":true,
		"withoutSubmission": true,
	}

	allowedPrice := map[string]bool{
		"0":true,
		"<100000":true,
		"<500000":true,
		"<1000000":true,
		"<2500000":true,
		">5000000":true,
	}

	for key := range c.QueryParams() {
		if !allowedParams[key] {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Query Parameter tidak dikenali"))
		}
	}

	filter := make(map[string]interface{})

	if price := c.QueryParam("price"); price != "" {
		if price != "" && !allowedPrice[price] {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Nilai Parameter price tidak dikenali"))
		}
		filter["price"] = price
	}

	if province := c.QueryParam("province"); province != "" {
		filter["province"] = province
	}

	if timeStart := c.QueryParam("timeStart"); timeStart != "" {
		if timeStart != "" && !allowedTime[timeStart] {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Nilai Parameter timeStart tidak dikenali"))
		}
		filter["timeStart"] = timeStart
	}

	if category := c.QueryParam("category"); category != "" {
		if category != "" && !allowedCategory[category] {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Nilai Parameter category tidak dikenali"))
		}
		filter["category"] = category
	}

	sort := c.QueryParam("sort")

	if sort != "" && !allowedSorts[sort] {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Nilai Parameter sort tidak dikenali"))
	}

	events, err := h.eventService.GetEvents(filter, sort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if len(events) == 0 {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Data Event Tidak Ditemukan Satupun", nil))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Berhasil Menampilkan Events",events))
}
