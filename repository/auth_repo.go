package repository

import (
	"satixnimble/identity/model"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type AuthRepo struct {
	logger *logrus.Entry
	db     *gorm.DB
}

func NewAuthRepo(logger *logrus.Entry, db *gorm.DB) *AuthRepo {
	return &AuthRepo{
		logger: logger,
		db:     db,
	}
}

func (r *AuthRepo) GetUser(username string) (*model.Users, error) {
	user := &model.Users{}
	if err := r.db.Where(&model.Users{
		Username: username,
	}).Find(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
