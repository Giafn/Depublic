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
	excludeUser         = []string{Admin, PetugasLapangan}
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
		{
			Method:  http.MethodGet,
			Path:    "/payment",
			Handler: transactionHandler.PaymentRedirect,
		},
	}
}

func AppPrivateRoutes(Handler handler.AppHandler) []*route.Route {
	userHandler := Handler.UserHandler
	profileHandler := Handler.ProfileHandler
	transactionHandler := Handler.TransactionHandler
	eventHandler := Handler.EventHandler
	ticketHandler := Handler.TicketHandler
	submissionHandler := Handler.SubmissionHandler
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
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/user",
			Handler: ticketHandler.FindTicketsByUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:id",
			Handler: ticketHandler.FindTicketByID,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/booking/:bookingNum",
			Handler: ticketHandler.FindTicketByBookingNumber,
			Roles:   excludeUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/:id",
			Handler: ticketHandler.UpdateTicket,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket/validate/:id",
			Handler: ticketHandler.ValidateTicket,
			Roles:   onlyPetugasLapangan,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ticket/:id",
			Handler: ticketHandler.DeleteTicketById,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/transaction/:transactionId",
			Handler: ticketHandler.FindTicketsByTransactionId,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ticket/booking/:bookingNum",
			Handler: ticketHandler.DeleteTicketByBookingNumber,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/mytransactions",
			Handler: transactionHandler.FindMyTransactions,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/mytransactions",
			Handler: transactionHandler.FindMyTransactions,
			Roles:   onlyUser,
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
			Roles:   allRoles,
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
			Handler: notificationHandler.MarkAllNotificationsAsSeen,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/notifications",
			Handler: notificationHandler.DeleteSeenNotifications,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/notifications/:id",
			Handler: notificationHandler.DeleteNotificationByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/submission",
			Handler: submissionHandler.CreateSubmission,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submission",
			Handler: submissionHandler.ListSubmission,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/submission/accept/:id",
			Handler: submissionHandler.AcceptSubmission,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/submission/reject/:id",
			Handler: submissionHandler.RejectSubmission,
			Roles:   allRoles,
		},
	}
}
