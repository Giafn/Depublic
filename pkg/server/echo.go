package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Giafn/Depublic/pkg/response"
	"github.com/Giafn/Depublic/pkg/route"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(serverName string, publicRoutes, privateRoutes []*route.Route, secretKey string) *Server {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Welcome To Depublic App please Access /app/api/v1/register for register!", nil))
	})

	v1 := e.Group(fmt.Sprintf("%s/api/v1", serverName))

	if len(publicRoutes) > 0 {
		for _, v := range publicRoutes {
			v1.Add(v.Method, v.Path, v.Handler)
		}
	}

	if len(privateRoutes) > 0 {
		for _, v := range privateRoutes {
			v1.Add(v.Method, v.Path, v.Handler, JWTProtection(secretKey), RBACMiddleware(v.Roles...))
		}
	}

	return &Server{e}
}

func (s *Server) Run(port string) {
	runServer(s, port)
	gracefulShutdown(s)
}

func runServer(srv *Server, port string) {
	go func() {
		err := srv.Start(fmt.Sprintf(":%s", port))
		log.Fatal(err)
	}()
}

func gracefulShutdown(srv *Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger.Fatal("Server Shutdown:", err)
		}
	}()
}

func JWTProtection(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(secretKey),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "anda harus login untuk mengakses resource ini"))
		},
	})
}

func RBACMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "anda harus login untuk mengakses resource ini"))
			}

			claims := user.Claims.(*token.JwtCustomClaims)

			// Check if the user has the required role
			if !contains(roles, claims.Role) {
				return c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "anda tidak memiliki akses untuk resource ini"))
			}

			return next(c)
		}
	}
}

func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
