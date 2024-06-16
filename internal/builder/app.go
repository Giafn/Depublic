package builder

import (
	"github.com/Giafn/Depublic/configs"
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

func BuildAppPublicRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, tokenUseCase token.TokenUseCase, cfg *configs.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)

	profileRepository := repository.NewProfileRepository(db, cacheable)

	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, profileRepository, tokenUseCase, encryptTool, cfg)
	userHandler := handler.NewUserHandler(userService)

	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService)

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository, db)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	appHandler := handler.NewAppHandler(userHandler, transactionHandler, ticketHandler, nil)
	return router.AppPublicRoutes(appHandler)
}

func BuildAppPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, cfg *configs.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)

	profileRepository := repository.NewProfileRepository(db, cacheable)
	profileService := service.NewProfileService(profileRepository, encryptTool)
	profileHandler := handler.NewProfileHandler(profileService)

	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, profileRepository, nil, encryptTool, cfg)
	userHandler := handler.NewUserHandler(userService)

	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService)

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository, db)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	appHandler := handler.NewAppHandler(userHandler, transactionHandler, ticketHandler, profileHandler)
	return router.AppPrivateRoutes(appHandler)
}
