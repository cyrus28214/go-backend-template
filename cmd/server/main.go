package main

import (
	"backend/internal/user"
	"backend/pkg/config"
	"backend/pkg/database"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"backend/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// config
	cfg := config.LoadConfig("configs")

	// logger
	logger, err := logger.InitLogger(&cfg.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	logger.Info("logger initialized")

	// database
	DB, err := database.InitDatabase(&cfg.Database, logger.With("module", "database"))
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	logger.Info("database initialized")

	err = DB.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	logger.Info("database migrated")

	// repository
	userRepo := user.NewUserRepository(DB)

	// service
	userService := user.NewUserService(userRepo, &cfg.Jwt)

	// handler
	userHandler := user.NewUserHandler(userService)

	// router
	router := gin.New()
	switch cfg.Server.Mode {
	case "development":
		gin.SetMode(gin.DebugMode)
	case "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		log.Fatalf("invalid server mode: %s", cfg.Server.Mode)
	}
	logger.Info("router initialized")

	// middleware
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware)

	// routes
	router.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})

	router.POST("/api/v1/auth/login", userHandler.Login)

	// server
	logger.Info("server will be started at " + cfg.Server.Address)
	if err := router.Run(cfg.Server.Address); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
