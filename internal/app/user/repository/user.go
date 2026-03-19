// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
)

type UserDBItf interface {
	Register(user *entity.User) error
	RegisterUserDetail(userDetail *entity.UserDetail) error
	UpdateUserDetail(userDetail *entity.UserDetail) error
	UpdateUserInfo(user *entity.User) error
	Login(user *entity.User) error
	GoogleOAuth(user *entity.User) error
	CheckUsername(user *entity.User) error
	GetUserIDFromUsername(user *entity.User) error
	GetUsername(user *entity.User, userParam dto.Login) error
	GetUserInfo(user *entity.User) error
	SoftDelete(user *entity.User) error
}

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) UserDBItf {
	return &UserDB{
		db: db,
	}
}

func (r *UserDB) Register(user *entity.User) error {
	return r.db.Debug().
		Create(user).
		Error
}

func (r *UserDB) RegisterUserDetail(userDetail *entity.UserDetail) error {
	return r.db.Debug().
		Create(userDetail).
		Error
}

func (r *UserDB) UpdateUserInfo(user *entity.User) error {
	return r.db.Debug().
		Updates(user).
		Error
}

func (r *UserDB) UpdateUserDetail(userDetail *entity.UserDetail) error {
	var err error

	if r.db.Debug().
		Model(&userDetail).
		Where("user_id = ?", userDetail.UserID).Error != nil {

		err = r.db.Create(userDetail).Error
	} else {
		err = r.db.Debug().
			Updates(&userDetail).
			Error
	}

	return err
}

func (r *UserDB) Login(user *entity.User) error {
	return r.db.Debug().
		First(user).
		Error
}

func (r *UserDB) GoogleOAuth(user *entity.User) error {
	// TODO: Fix

	r.db.Debug().
		Create(user)

	return r.db.Debug().
		Model(user).
		Where("users.email = ?", user.Email).
		First(user).
		Error
}

func (r *UserDB) CheckUsername(user *entity.User) error {
	return r.db.Debug().
		Raw("SELECT `username` FROM `users` WHERE username = ?", user.Username).
		First(&user).
		Error
}

func (r *UserDB) GetUsername(user *entity.User, userParam dto.Login) error {
	return r.db.Debug().
		First(user, userParam).
		Error
}

func (r *UserDB) GetUserInfo(user *entity.User) error {
	return r.db.Debug().
		Model(&user).
		Preload("UserDetail").
		Select("users.id, users.email, users.username, users.name, users.created_at, users.updated_at, user_details.*").
		Joins("LEFT JOIN user_details ON user_details.user_id = users.id").
		First(&user).
		Error
}

func (r *UserDB) GetUserIDFromUsername(user *entity.User) error {
	return r.db.Debug().
		Select("id").
		Where("username = ?", user.Username).
		First(user).
		Error
}

func (r *UserDB) SoftDelete(user *entity.User) error {
	return r.db.Debug().
		Delete(user).
		Error
}
