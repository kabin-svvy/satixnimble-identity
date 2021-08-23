package router

import (
	"satixnimble/identity/database"
	"satixnimble/identity/handler"
	"satixnimble/identity/middleware"
	"satixnimble/identity/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {

	log.SetFormatter(&log.JSONFormatter{})

	api := app.Group("/api", logger.New())
	api.Get("/", handler.HealthCheck)

	authRepo := repository.NewAuthRepo(log.WithField("Repository", "AuthRepo"), database.DB)
	authHandle := handler.NewAuthHandle(log.WithField("Handler", "AuthHandler"), authRepo)
	auth := api.Group("/auth")
	auth.Post("/login", authHandle.Login)

	userRepo := repository.NewUserRepo(log.WithField("Repository", "UserRepo"), database.DB)
	userHandle := handler.NewUserHandle(log.WithField("Handler", "UserHandler"), userRepo)
	user := api.Group("/user")
	user.Get("/:id", userHandle.GetUser)
	user.Post("/", userHandle.CreateUser)
	user.Patch("/:id", middleware.Protected(), userHandle.UpdateUser)
	user.Delete("/:id", middleware.Protected(), userHandle.DeleteUser)
}
