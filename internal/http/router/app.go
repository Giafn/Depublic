package router

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/http/handler"
	"github.com/Giafn/Depublic/pkg/route"
)

const (
	Admin           = "Admin"
	User            = "User"
	PetugasLapangan = "PetugasLapangan"
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
	transactionHandler := Handler.TransactionHandler
	eventHandler := Handler.EventHandler

	return []*route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: welcome,
		},
		{
			Method:  http.MethodGet,
			Path:    "/file/:filepath",
			Handler: Handler.FileReader,
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
		{
			Method:  http.MethodPost,
			Path:    "/transaction/webhook",
			Handler: transactionHandler.WebhookPayment,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event",
			Handler: eventHandler.GetEvents,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/:id",
			Handler: eventHandler.FindEventByID,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/pricing/:id",
			Handler: eventHandler.FindPricingByEventID,
		},
	}
}

func AppPrivateRoutes(Handler handler.AppHandler) []*route.Route {
	userHandler := Handler.UserHandler
	profileHandler := Handler.ProfileHandler
	transactionHandler := Handler.TransactionHandler
	eventHandler := Handler.EventHandler
	ticketHandler := Handler.TicketHandler
	notificationHandler := Handler.NotificationHandler

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
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.FindUserByID,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/:id",
			Handler: userHandler.UpdateUser,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/profile",
			Handler: profileHandler.FindCurrentUserProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/profile",
			Handler: profileHandler.UpdateProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/profile",
			Handler: profileHandler.DeleteProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/event",
			Handler: eventHandler.CreateNewEvent,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event/:id",
			Handler: eventHandler.UpdateEventWithPricing,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event/:id",
			Handler: eventHandler.DeleteEvent,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event/pricing/:id",
			Handler: eventHandler.DeletePricing,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/event/pricing",
			Handler: eventHandler.CreatePricing,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/create",
			Handler: ticketHandler.CreateTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket",
			Handler: ticketHandler.FindAllTickets,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:id",
			Handler: ticketHandler.FindTicketByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:bookingNum",
			Handler: ticketHandler.FindTicketByBookingNumber,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/:id",
			Handler: ticketHandler.UpdateTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/:id/validate",
			Handler: ticketHandler.ValidateTicket,
			Roles:   onlyPetugasLapangan,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ticket/:id",
			Handler: ticketHandler.DeleteTicketById,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ticket/:bookingNum",
			Handler: ticketHandler.DeleteTicketByBookingNumber,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/transactions",
			Handler: transactionHandler.CreateTransaction,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/:id",
			Handler: transactionHandler.FindTransactionByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions",
			Handler: transactionHandler.FindAllTransactions,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/transactions/:id",
			Handler: transactionHandler.UpdateTransaction,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/transactions/:id",
			Handler: transactionHandler.DeleteTransaction,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/notifications",
			Handler: notificationHandler.FindAllNotification,
			Roles:  allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/notifications/:id",
			Handler: notificationHandler.FindNotificationByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/notifications",
			Handler: notificationHandler.UpdateSeenNotifications,
			Roles:  allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/notifications",
			Handler: notificationHandler.DeleteSeenNotifications,
			Roles:  allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/notifications/:id",
			Handler: notificationHandler.DeleteNotificationByID,
			Roles:  allRoles,
		},
		
	}
}
