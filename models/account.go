package models

import (
	"gorm.io/gorm"
	"weibo/pkg/mysql"
)

type Account struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Role     string
	Salt     string
	Image    string
}

func ExistAccountByEmail(email string) (bool, error) {
	var account Account
	result := mysql.DB.Where(Account{Email: email}).First(&account)
	err := result.Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func FindAccountByEmail(email string) (*Account, error) {
	var account Account
	result := mysql.DB.Where(Account{Email: email}).First(&account)
	err := result.Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if result.RowsAffected > 0 {
		return &account, nil
	}

	return nil, nil
}

func CreateUser(account Account) *gorm.DB {
	return mysql.DB.Create(&account)
}
