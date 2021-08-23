package repository

import (
	"fmt"
	"satixnimble/identity/model"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type UserRepo struct {
	logger *logrus.Entry
	db     *gorm.DB
}

func NewUserRepo(logger *logrus.Entry, db *gorm.DB) *UserRepo {
	return &UserRepo{
		logger: logger,
		db:     db,
	}
}

func (r *UserRepo) CreateUser(users *model.Users) error {
	if err := r.db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UpdateUser(id string, firstname string) (*model.Users, error) {

	var users model.Users
	r.db.First(&users, id)

	users.Firstname = firstname

	if err := r.db.Save(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *UserRepo) GetUser(id string) (*model.Users, error) {
	user := &model.Users{}
	r.db.Find(user, id)
	if user.Username == "" {
		err := fmt.Errorf("%v", "User not found with id")
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) DeleteUser(id string) (string, error) {
	user := &model.Users{}
	r.db.First(user, id)
	if user.Username == "" {
		err := fmt.Errorf("%v", "User not found with id")
		return "", err
	}
	username := user.Username
	r.db.Delete(user)
	return username, nil
}
