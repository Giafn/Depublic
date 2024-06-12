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
	profileHandler := Handler.ProfileHandler

	return []*route.Route{
	
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.FindUserByID,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAllUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/profile",
			Handler: profileHandler.CreateProfile,
		},
		{
			Method:  http.MethodGet,
			Path:    "/profile",
			Handler: profileHandler.FindCurrentUserProfile,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/profile",
			Handler: profileHandler.UpdateProfile,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/profile",
			Handler: profileHandler.DeleteProfile,
		},
	}
}
