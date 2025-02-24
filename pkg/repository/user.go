package repository

import (
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/tools"
	"gorm.io/gorm"
)

type UserRepo struct {
	mdb *gorm.DB
}

var userRepo *UserRepo

func NewUserRepo(db *gorm.DB) *UserRepo {
	if db == nil {
		tools.ErrorMsg("DB connection is nil")
		return nil
	}
	tools.Log("New user repo ...")
	u := &UserRepo{
		mdb: db,
	}
	userRepo = u
	return u
}

func GetUserRepo() *UserRepo {
	return userRepo
}

func (userRepo *UserRepo) CreateUser(user mdb.UserModel) {
	var count int64
	userRepo.mdb.Model(&mdb.UserModel{}).Where("email = ?", user.Email).Count(&count)
	if count < 0 {
		tx := userRepo.mdb.Begin()
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			tools.HandelError("error in create user", err)
		}
		tx.Commit()
	} else {
		tools.ErrorMsg("have the same user email ")
	}
}
