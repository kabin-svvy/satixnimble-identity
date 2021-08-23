package handler

import (
	"satixnimble/identity/common"
	"satixnimble/identity/middleware"
	"satixnimble/identity/model"
	"satixnimble/identity/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler interface {
	CreateUser(*fiber.Ctx) error
	UpdateUser(*fiber.Ctx) error
	GetUser(*fiber.Ctx) error
}

type UserHandle struct {
	logger *logrus.Entry
	repo   *repository.UserRepo
}

func NewUserHandle(logger *logrus.Entry, repo *repository.UserRepo) *UserHandle {
	return &UserHandle{
		logger: logger,
		repo:   repo,
	}
}

func (r *UserHandle) CreateUser(c *fiber.Ctx) error {
	type NewUserInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	user := new(model.Users)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err.Error()})
	}

	hash, err := common.HashPassword(user.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": err.Error()})
	}

	user.Password = hash

	if err := r.repo.CreateUser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err.Error()})
	}

	newUser := NewUserInput{
		Username: user.Username,
		Email:    user.Email,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "create user", "data": newUser})
}

func (r *UserHandle) UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Firstname string `json:"firstname"`
	}

	type UpdateUserOutPut struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
	}

	var updateUserInput UpdateUserInput

	if err := c.BodyParser(&updateUserInput); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err.Error()})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !middleware.ValidToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	user, err := r.repo.UpdateUser(id, updateUserInput.Firstname)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't update user", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "update user", "data": UpdateUserOutPut{
		Username:  user.Username,
		Firstname: user.Firstname,
	}})
}

func (r *UserHandle) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := r.repo.GetUser(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "get user", "data": user})
}

func (r *UserHandle) DeleteUser(c *fiber.Ctx) error {
	type DeleteUserInput struct {
		Password string `json:"password"`
	}

	type DeleteUserOutput struct {
		Username string `json:"username"`
	}

	var deleteUserInput DeleteUserInput

	if err := c.BodyParser(&deleteUserInput); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err.Error()})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !middleware.ValidToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	user, err := r.repo.GetUser(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid information", "data": err.Error()})
	}

	if !common.CheckPasswordHash(deleteUserInput.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": nil})
	}

	username, err := r.repo.DeleteUser(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't delete user", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "delete user", "data": DeleteUserOutput{
		Username: username,
	}})
}
