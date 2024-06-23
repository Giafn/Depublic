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
	profileService := service.NewProfileService(profileRepository, encryptTool)
	profileHandler := handler.NewProfileHandler(profileService)

	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, profileRepository, tokenUseCase, encryptTool, cfg)
	userHandler := handler.NewUserHandler(userService)

	eventRepository := repository.NewEventRepository(db)
	pricingRepository := repository.NewPricingRepository(db)
	eventService := service.NewEventService(eventRepository, pricingRepository)
	eventHandler := handler.NewEventHandler(eventService)

	notificationRepository := repository.NewNotificationRepository(db, cacheable)
	notificationService := service.NewNotificationService(notificationRepository)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	ticketRepository := repository.NewTicketRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(
		transactionRepository,
		pricingRepository,
		userRepository,
		eventRepository,
		ticketRepository,
		notificationRepository,
		db,
		encryptTool,
		cfg,
	)

	ticketHandler := handler.NewTicketHandler(ticketService, transactionService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	appHandler := handler.NewAppHandler(userHandler,
		transactionHandler,
		ticketHandler,
		profileHandler,
		eventHandler,
		notificationHandler,
		nil,
	)
	return router.AppPublicRoutes(appHandler)
}

func BuildAppPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, tokenUseCase token.TokenUseCase, cfg *configs.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)

	profileRepository := repository.NewProfileRepository(db, cacheable)
	profileService := service.NewProfileService(profileRepository, encryptTool)
	profileHandler := handler.NewProfileHandler(profileService)

	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, profileRepository, tokenUseCase, encryptTool, cfg)
	userHandler := handler.NewUserHandler(userService)

	eventRepository := repository.NewEventRepository(db)
	pricingRepository := repository.NewPricingRepository(db)
	eventService := service.NewEventService(eventRepository, pricingRepository)
	eventHandler := handler.NewEventHandler(eventService)

	notificationRepository := repository.NewNotificationRepository(db, cacheable)
	notificationService := service.NewNotificationService(notificationRepository)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	ticketRepository := repository.NewTicketRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(
		transactionRepository,
		pricingRepository,
		userRepository,
		eventRepository,
		ticketRepository,
		notificationRepository,
		db,
		encryptTool,
		cfg,
	)

	ticketHandler := handler.NewTicketHandler(ticketService, transactionService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	submissionRepository := repository.NewSubmissionRepository(db, cacheable)
	submissionService := service.NewSubmissionService(submissionRepository, transactionRepository, userRepository, eventRepository, notificationRepository, cfg)
	submissionHandler := handler.NewSubmissionHandler(submissionService)

	appHandler := handler.NewAppHandler(userHandler,
		transactionHandler,
		ticketHandler,
		profileHandler,
		eventHandler,
		notificationHandler,
		submissionHandler,
	)
	return router.AppPrivateRoutes(appHandler)
}
