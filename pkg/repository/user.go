package repository

import (
	"errors"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/utils"
	"gorm.io/gorm"
)

type UserRepo struct {
	mdb *gorm.DB
}

var userRepo *UserRepo

func NewUserRepo(db *gorm.DB) *UserRepo {
	if db == nil {
		utils.ErrorMsg("DB connection is nil")
		return nil
	}
	u := &UserRepo{
		mdb: db,
	}
	userRepo = u
	return u
}

func GetUserRepo() *UserRepo {
	return userRepo
}

func (userRepo *UserRepo) CreateUser(user mdb.UserModel) bool {
	var count int64
	userRepo.mdb.Model(&mdb.UserModel{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		utils.ErrorMsg("user already exists")
		return false
	}
	tx := userRepo.mdb.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		utils.HandelError("error in create user", err)
	}
	tx.Commit()
	return true
}

func (userRepo *UserRepo) LoginUserByEmailAndPassword(user mdb.UserModel) (int64, string) {
	res := userRepo.mdb.Where("email = ? ", user.Email).Where("password =?", user.Password).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		utils.ErrorMsg("not have the user")
		return -1, ""
	} else if res.Error != nil {
		utils.HandelError("error in loginUserByEmailAndPassword", res.Error)
		return -1, ""
	}
	return user.Id, user.Name
}
