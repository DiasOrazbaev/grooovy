package main

import (
	"grovo/config"
	"grovo/internal/controller/http"
	"grovo/internal/usecase"
	"grovo/internal/usecase/repo"
	mongo2 "grovo/pkg/mongo"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	var (
		cfg *config.Config
		db  *mongo.Database
	)

	app := fiber.New(fiber.Config{
		AppName:     "groovy-backend",
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal("config", zap.Error(err))
	}

	db, err = mongo2.NewMongoDatabase(cfg)
	if err != nil {
		log.Fatal("mongo", zap.Error(err))
	}

	// repo
	userRepo := repo.NewMongoUserRepository(db)
	// service
	userUsecase := usecase.NewUserUsecase(userRepo)
	// handler
	userController := http.NewUserController(userUsecase, log)
	userController.Routes(app)

	// graceful shutdown
	go func() {
		if err := app.Listen(":3005"); err != nil {
			log.Fatal("server shutdown", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if !fiber.IsChild() {
		log.Info("shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatal("server shutdown", zap.Error(err))
		}

		log.Info("server gracefully stopped")
	}
}
