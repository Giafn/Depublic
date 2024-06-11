package router

import (
	"net/http"

	"github.com/Giafn/Depublic/internal/http/handler"
	"github.com/Giafn/Depublic/pkg/route"
)

const (
	Admin = "Admin"
	User  = "User"
)

var (
	allRoles  = []string{Admin, User}
	onlyAdmin = []string{Admin}
	onlyUser  = []string{User}
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

	return []*route.Route{
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
	}
}
