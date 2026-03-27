// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/app/user/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/SyafaHadyan/worku/internal/infra/jwt"
	redisitf "github.com/SyafaHadyan/worku/internal/infra/redis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCaseItf interface {
	Register(register dto.Register) (dto.ResponseRegister, error)
	UpdateUserInfo(updateUserInfo dto.UpdateUserInfo, userID uuid.UUID) (dto.ResponseUpdateUserInfo, error)
	Login(login dto.Login) (dto.ResponseLogin, string, error)
	GoogleOAuth(responseGoogleOAuth dto.ResponseGoogleOAuth) (dto.ResponseLogin, string, error)
	GetUserIDFromUsername(username string) (uuid.UUID, error)
	GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error)
	SoftDelete(userID uuid.UUID) error
}

type UserUseCase struct {
	userRepo     repository.UserDBItf
	jwt          jwt.JWTItf
	redis        redisitf.RedisItf
	redisContext context.Context
}

func NewUserUseCase(
	userRepo repository.UserDBItf, jwt *jwt.JWT,
	redis redisitf.RedisItf,
) UserUseCaseItf {
	return &UserUseCase{
		userRepo:     userRepo,
		jwt:          jwt,
		redis:        redis,
		redisContext: context.Background(),
	}
}

func (u *UserUseCase) Register(register dto.Register) (dto.ResponseRegister, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(register.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	user := entity.User{
		ID:       uuid.New(),
		Email:    register.Email,
		Username: register.Username,
		Password: string(hashedPassword),
		Name:     register.Name,
	}

	userDetail := entity.UserDetail{
		UserID: user.ID,
	}

	err = u.userRepo.Register(&user)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	err = u.userRepo.RegisterUserDetail(&userDetail)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	user.UserDetail = userDetail

	return user.ParseToDTOResponseRegister(), nil
}

func (u *UserUseCase) UpdateUserInfo(updateUserInfo dto.UpdateUserInfo, userID uuid.UUID) (dto.ResponseUpdateUserInfo, error) {
	var redisKey string

	user := entity.User{
		ID:       userID,
		Email:    updateUserInfo.Email,
		Username: updateUserInfo.Username,
		Name:     updateUserInfo.Name,
	}

	userDetail := entity.UserDetail{
		UserID: userID,
	}

	err := u.userRepo.UpdateUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	err = u.userRepo.UpdateUserDetail(&userDetail)
	if err != nil {
		log.Println(err)
	}

	err = u.userRepo.GetUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	redisKey = fmt.Sprintf("user:%s", userID.String())

	go func() {
		newData, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		u.redis.Set(redisKey, string(newData))
	}()

	return user.ParseToDTOResponseUpdateUserInfo(), nil
}

func (u *UserUseCase) Login(login dto.Login) (dto.ResponseLogin, string, error) {
	user := entity.User{
		Username: login.Username,
		Email:    login.Email,
	}

	err := u.userRepo.Login(&user)
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	_ = u.userRepo.GetUserInfo(&user)

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	return user.ParseToDTOResponseLogin(), token, nil
}

func (u *UserUseCase) GoogleOAuth(responseGoogleOAuth dto.ResponseGoogleOAuth) (dto.ResponseLogin, string, error) {
	user := entity.User{
		ID:       uuid.New(),
		Email:    responseGoogleOAuth.Email,
		Username: responseGoogleOAuth.Email,
		Name:     responseGoogleOAuth.Email,
	}

	err := u.userRepo.GoogleOAuth(&user)
	if err != gorm.ErrRecordNotFound {
		return dto.ResponseLogin{},
			"",
			err
	}

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	return user.ParseToDTOResponseLogin(), token, nil
}

func (u *UserUseCase) CheckUsername(userName *dto.CheckUsername) error {
	user := entity.User{
		Username: userName.Username,
	}

	err := u.userRepo.CheckUsername(&user)

	return err
}

func (u *UserUseCase) GetUserIDFromUsername(username string) (uuid.UUID, error) {
	user := entity.User{
		Username: username,
	}

	key := fmt.Sprintf("user:%s", username)

	result, err := u.redis.Get(key)
	if err == nil && result != "" {
		userID, _ := uuid.Parse(result)

		return userID, nil
	}

	err = u.userRepo.GetUserIDFromUsername(&user)
	if err != nil {
		return uuid.Nil,
			err
	}

	go func() {
		u.redis.Set(key, user.ID.String())
	}()

	return user.ID, nil
}

func (u *UserUseCase) GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error) {
	user := entity.User{
		ID: userID,
	}

	redisKey := fmt.Sprintf("user:%s", userID.String())

	result, err := u.redis.Get(redisKey)
	if err == nil && result != "" {
		var out entity.User

		err := json.Unmarshal([]byte(result), &out)
		if err != nil {
			log.Println(err)
		}

		return out.ParseToDTOResponseGetUserInfo(), nil
	}

	err = u.userRepo.GetUserInfo(&user)
	if err != nil {
		return dto.ResponseGetUserInfo{},
			err
	}

	go func() {
		newData, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		u.redis.Set(redisKey, string(newData))
	}()

	return user.ParseToDTOResponseGetUserInfo(), nil
}

func (u *UserUseCase) SoftDelete(userID uuid.UUID) error {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.SoftDelete(&user)

	return err
}
