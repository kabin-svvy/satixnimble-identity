package handler

import (
	"satixnimble/identity/common"
	"satixnimble/identity/config"
	"satixnimble/identity/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthHandler interface {
	Login(*fiber.Ctx) error
}

type AuthHandle struct {
	logger *logrus.Entry
	repo   *repository.AuthRepo
}

func NewAuthHandle(logger *logrus.Entry, repo *repository.AuthRepo) *AuthHandle {
	return &AuthHandle{
		logger: logger,
		repo:   repo,
	}
}

func (r *AuthHandle) Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err})
	}
	username := input.Username
	password := input.Password

	user, err := r.repo.GetUser(username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": err})
	}

	if !common.CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
