// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Email      string         `json:"email" gorm:"type:nvarchar(256);not null;unique"`
	Username   string         `json:"username" gorm:"type:nvarchar(64);not null;unique"`
	Password   string         `json:"password" gorm:"type:text;not null"`
	Name       string         `json:"name" gorm:"type:nvarchar(128)"`
	CreatedAt  time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserDetail UserDetail
}

type UserDetail struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
}

func (u *User) ParseToDTOResponseRegister() dto.ResponseRegister {
	var responseRegister dto.ResponseRegister

	responseRegister.ID = u.ID
	responseRegister.Email = u.Email
	responseRegister.Username = u.Username
	responseRegister.Name = u.Name
	responseRegister.CreatedAt = u.CreatedAt
	responseRegister.UpdatedAt = u.UpdatedAt

	return responseRegister
}

func (u *User) ParseToDTOResponseLogin() dto.ResponseLogin {
	var responseLogin dto.ResponseLogin

	responseLogin.ID = u.ID
	responseLogin.Email = u.Email
	responseLogin.Username = u.Username
	responseLogin.Name = u.Name
	responseLogin.CreatedAt = u.CreatedAt
	responseLogin.UpdatedAt = u.UpdatedAt

	return responseLogin
}

func (u *User) ParseToDTOResponseGetUserInfo() dto.ResponseGetUserInfo {
	var responseGetUserInfo dto.ResponseGetUserInfo

	responseGetUserInfo.ID = u.ID
	responseGetUserInfo.Email = u.Email
	responseGetUserInfo.Username = u.Username
	responseGetUserInfo.Name = u.Name
	responseGetUserInfo.CreatedAt = u.CreatedAt
	responseGetUserInfo.UpdatedAt = u.UpdatedAt

	return responseGetUserInfo
}

func (u *User) ParseToDTOResponseGetUserInfoPublic() dto.ResponseGetUserInfoPublic {
	var responseGetUserInfoPublic dto.ResponseGetUserInfoPublic

	responseGetUserInfoPublic.Username = u.Username
	responseGetUserInfoPublic.Name = u.Name

	return responseGetUserInfoPublic
}

func (u *User) ParseToDTOResponseUpdateUserInfo() dto.ResponseUpdateUserInfo {
	var responseUdpateUserInfo dto.ResponseUpdateUserInfo

	responseUdpateUserInfo.ID = u.ID
	responseUdpateUserInfo.Email = u.Email
	responseUdpateUserInfo.Username = u.Username
	responseUdpateUserInfo.Name = u.Name
	responseUdpateUserInfo.CreatedAt = u.CreatedAt
	responseUdpateUserInfo.UpdatedAt = u.UpdatedAt

	return responseUdpateUserInfo
}
