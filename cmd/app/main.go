package main

import (
	"time"

	"github.com/Giafn/Depublic/configs"
	"github.com/Giafn/Depublic/internal/builder"
	"github.com/Giafn/Depublic/pkg/background_job"
	"github.com/Giafn/Depublic/pkg/cache"
	"github.com/Giafn/Depublic/pkg/encrypt"
	"github.com/Giafn/Depublic/pkg/postgres"
	"github.com/Giafn/Depublic/pkg/server"
	"github.com/Giafn/Depublic/pkg/token"
)

func main() {
	cfg, err := configs.NewConfig()
	checkError(err)

	db, err := postgres.InitDB(&cfg.Postgres)
	checkError(err)

	redisDB, err := cache.InitRedis(&cfg.Redis)
	checkError(err)

	go background_job.Init(cfg.Redis.Host, cfg.Redis.Port)

	tokenUse := token.NewTokenUseCase(cfg.JWT.SecretKey, time.Duration(cfg.JWT.ExpiresAt)*time.Hour)
	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.Iv)

	publicRoutes := builder.BuildAppPublicRoutes(db, redisDB, encryptTool, tokenUse, cfg)
	privateRoutes := builder.BuildAppPrivateRoutes(db, redisDB, encryptTool, cfg)

	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey, tokenUse)
	srv.Run(cfg.Port)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
