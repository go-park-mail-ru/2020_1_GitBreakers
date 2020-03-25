package repository

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DBWork struct {
	Db *sqlx.DB
}

func (db *DBWork) Create(user *models.User) error {
	//insert to db User.email, user.password, user.nickname
	return nil
}

func (db *DBWork) Update(user *models.User) error {
	//update to db User.email, user.password, user.nickname
	return nil
}

func (db *DBWork) SaveAvatar(user *models.User, filepath string) error {
	//update to db User.email, user.password, user.nickname
	return nil
}
