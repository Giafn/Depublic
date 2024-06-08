package handler

import (
	"net/http"

	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AppHandler struct {
	WelcomeHandler     echo.HandlerFunc
	UserHandler        UserHandler
	TransactionHandler TransactionHandler
}

func NewAppHandler(userHandler UserHandler, transactionHandler TransactionHandler) AppHandler {
	return AppHandler{
		WelcomeHandler:     welcome,
		UserHandler:        userHandler,
		TransactionHandler: transactionHandler,
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