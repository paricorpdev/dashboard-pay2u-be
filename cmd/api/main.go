package main

import (
	"fmt"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	redisStorage := config.NewRedis(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Redis:    redisStorage,
	})

	timeZone := viperConfig.GetString("app.timezone")
	if loc, err := time.LoadLocation(timeZone); err == nil {
		time.Local = loc
	}

	// graceful shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		webPort := viperConfig.GetInt("web.port")
		err := app.Listen(fmt.Sprintf(":%d", webPort))
		if err != nil {
			log.Fatalf("failed to start server: %+v", err)
		}
	}()

	<-shutdown
	log.Info("Gracefully shutting down...")

	defer func(db *gorm.DB) {
		connection, err := db.DB()
		if err != nil {
			log.Fatalf("failed to close database: %+v", err)
		}
		err = connection.Close()
		if err != nil {
			log.Fatalf("failed to close database: %+v", err)
		}
	}(db)

	defer func(app *fiber.App) {
		err := app.Shutdown()
		if err != nil {
			log.Fatalf("failed to shutdown server: %+v", err)
		}
	}(app)

	defer func() {
		if redisStorage != nil {
			_ = redisStorage.Close()
		}
	}()

	log.Info("Server stopped")
}
