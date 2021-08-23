package server

import (
	"fmt"
	"os"
	"os/signal"
	"satixnimble/identity/config"
	"satixnimble/identity/database"
	"satixnimble/identity/router"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

const idleTimeout = 5 * time.Second

func Run() {
	log.SetFormatter(&log.JSONFormatter{})

	port := config.Config("SERVER_PORT")

	database.ConnectDB()
	defer database.DB.Close()

	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	app.Use(cors.New())

	router.SetupRoutes(app)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", port)); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	defer fmt.Println("Identity.API was successful shutdown.")
}
