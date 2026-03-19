// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"required,email"`
	Username string    `json:"username" validate:"required,min=3,max=32"`
	Password string    `json:"password" validate:"required,min=4"`
	Name     string    `json:"name" validate:"omitempty,min=3,max=64"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=4"`
}

type ResponseGoogleOAuth struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type RenewToken struct {
	ID uuid.UUID `json:"id"`
}

type UserDetail struct {
	UserID uuid.UUID `json:"user_id"`
}

type CheckUserID struct {
	ID uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
}

type GetUserID struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}

type GetUserInfoPublic struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}

type UpdateUserInfo struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=3,max=32"`
	Password string `json:"password" validate:"omitempty,min=4"`
	Name     string `json:"name" validate:"omitempty,min=3,max=64"`
}

type UpdateUserDetail struct {
	LastEducation string `json:"last_education" validate:"omitempty,oneof=SMA SMK D1 D2 D3 D4 S1 S2 S3"`
	Location      string `json:"location"`
}

type EmailVerification struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email" validate:"required,email"`
	Code  uint      `json:"code"`
}

type ValidateEmail struct {
	Email string `json:"email" validate:"required,email"`
	Code  uint   `json:"code" validate:"required,min=8"`
}

type CheckUsername struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
}

type ResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordWithID struct {
	ID uuid.UUID `json:"id" validate:"required,min=36,max=36"`
}

type CheckPasswordResetCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  uint   `json:"code" validate:"required,min=8"`
}

type ResetPasswordWithCode struct {
	Email            string    `json:"email" validate:"required,email"`
	Code             uint      `json:"code" validate:"required"`
	Password         string    `json:"password" validate:"required,min=4"`
	PasswordChangeID uuid.UUID `json:"password_change_id"`
}

type ChangePassword struct {
	Password string `json:"password" validate:"required,min=4"`
}

type ResponseRegister struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseLogin struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetUserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetUserInfoPublic struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type ResponseUpdateUserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResposneGetUserSkill struct {
	Skill string `json:"string"`
}
