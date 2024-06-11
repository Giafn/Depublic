package builder

import (
	"github.com/Giafn/Depublic/internal/http/handler"
	"github.com/Giafn/Depublic/internal/http/router"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/internal/service"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/Giafn/Depublic/pkg/encrypt"
	"github.com/Giafn/Depublic/pkg/route"
	"github.com/Giafn/Depublic/pkg/token"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildAppPublicRoutes(db *gorm.DB, tokenUseCase token.TokenUseCase) []*route.Route {
	userRepository := repository.NewUserRepository(db, nil)
	userService := service.NewUserService(userRepository, tokenUseCase)
	userHandler := handler.NewUserHandler(userService)

	appHandler := handler.NewAppHandler(userHandler)
	return router.AppPublicRoutes(appHandler)
}

func BuildAppPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, nil)
	userHandler := handler.NewUserHandler(userService)

	appHandler := handler.NewAppHandler(userHandler)
	return router.AppPrivateRoutes(appHandler)
}
