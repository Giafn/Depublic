package router

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/http/handler"
	"github.com/Giafn/Depublic/pkg/route"
)

const (
	Admin           = "Admin"
	User            = "User"
	PetugasLapangan = "Petugas Lapangan"
)

var (
	allRoles            = []string{Admin, User, PetugasLapangan}
	onlyAdmin           = []string{Admin}
	onlyUser            = []string{User}
	onlyPetugasLapangan = []string{PetugasLapangan}
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
		{
			Method:  http.MethodGet,
			Path:    "/account/verify/:id",
			Handler: userHandler.VerifyEmail,
		},
		{
			Method:  http.MethodPost,
			Path:    "/account/resend-email-verification",
			Handler: userHandler.ResendEmailVerification,
		},
	}
}

func AppPrivateRoutes(Handler handler.AppHandler) []*route.Route {
	userHandler := Handler.UserHandler
	transactionHandler := Handler.TransactionHandler
	ticketHandler := Handler.TicketHandler

	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: userHandler.Logout,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: userHandler.CreateUser,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAllUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.FindUserByID,
			Roles:   onlyAdmin,
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
		{
			Method:  http.MethodPost,
			Path:    "/transactions",
			Handler: transactionHandler.CreateTransaction,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/:id",
			Handler: transactionHandler.FindTransactionByID,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions",
			Handler: transactionHandler.FindAllTransactions,
		},
		{
			Method:  http.MethodPut,
			Path:    "/transactions/:id",
			Handler: transactionHandler.UpdateTransaction,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/transactions/:id",
			Handler: transactionHandler.DeleteTransaction,
		},
	}
}
