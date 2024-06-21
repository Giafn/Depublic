package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AppHandler struct {
	WelcomeHandler     echo.HandlerFunc
	FileReader         echo.HandlerFunc
	UserHandler        UserHandler
	ProfileHandler     ProfileHandler
	EventHandler       EventHandler
	TransactionHandler TransactionHandler
	TicketHandler      TicketHandler
	NotificationHandler NotificationHandler
}

func NewAppHandler(userHandler UserHandler, transactionHandler TransactionHandler, ticketHandler TicketHandler, proprofileHandler ProfileHandler, eventHandler EventHandler, notificationHandler NotificationHandler) AppHandler {
	return AppHandler{
		WelcomeHandler:     welcome,
		FileReader:         fileReader,
		UserHandler:        userHandler,
		EventHandler:       eventHandler,
		TransactionHandler: transactionHandler,
		TicketHandler:      ticketHandler,
		ProfileHandler:     proprofileHandler,
		NotificationHandler: notificationHandler,
	}
}

func welcome(c echo.Context) error {
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "App Ticketing", nil))
}

func checkValidation(input interface{}) (errorMessage string, data interface{}) {
	validationErrors := validator.Validate(input)
	if validationErrors != nil {
		if _, exists := validationErrors["error"]; exists {
			return "validasi input gagal", nil
		}
		return "validasi input gagal", validationErrors
	}
	return "", nil
}

func fileReader(c echo.Context) error {
	type File struct {
		FilePath string `param:"filepath" validate:"required"`
	}
	input := new(File)

	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if errorMessage, data := checkValidation(input); errorMessage != "" {
		return c.JSON(http.StatusBadRequest, response.SuccessResponse(http.StatusBadRequest, errorMessage, data))
	}

	if input.FilePath[:7] != "uploads" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Not Found"))
	}

	return c.File(input.FilePath)
}
