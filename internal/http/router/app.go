package router

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/http/handler"
	"github.com/Giafn/Depublic/pkg/route"
)

func AppPublicRoutes(Handler handler.AppHandler) []*route.Route {
	welcome := Handler.WelcomeHandler
	userHandler := Handler.UserHandler
	return []*route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: welcome,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
	}
}

func AppPrivateRoutes(Handler handler.AppHandler) []*route.Route {
	userHandler := Handler.UserHandler
	ticketHandler := Handler.TicketHandler

	return []*route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAllUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.FindUserByID,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/create",
			Handler: ticketHandler.CreateTicket,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:id",
			Handler: ticketHandler.FindTicketByID,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/:id",
			Handler: ticketHandler.UpdateTicket,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/:id/validate",
			Handler: ticketHandler.ValidateTicket,
		},
		// delete
	}
}
