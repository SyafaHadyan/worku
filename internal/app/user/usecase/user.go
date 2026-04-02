// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/SyafaHadyan/worku/internal/app/user/repository"
	"github.com/SyafaHadyan/worku/internal/constants"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/SyafaHadyan/worku/internal/infra/jwt"
	redisitf "github.com/SyafaHadyan/worku/internal/infra/redis"
	s3itf "github.com/SyafaHadyan/worku/internal/infra/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCaseItf interface {
	Register(register dto.Register) (dto.ResponseRegister, error)
	Login(login dto.Login) (dto.ResponseLogin, string, error)
	GoogleOAuth(responseGoogleOAuth dto.ResponseGoogleOAuth) (dto.ResponseLogin, string, error)
	UploadProfilePicture(userID uuid.UUID, file multipart.FileHeader) (dto.ResponseGetUserInfo, error)
	GetUserIDFromUsername(username string) (uuid.UUID, error)
	GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error)
	GetUserDetail(userID uuid.UUID) (dto.ResponseGetUserDetail, error)
	GetUserContact(userID uuid.UUID) (dto.ResponseGetUserContact, error)
	GetUserEducation(userID uuid.UUID) (dto.ResponseGetUserEducation, error)
	GetUserLanguage(userID uuid.UUID) ([]dto.ResponseGetUserLanguage, error)
	GetUserEmployment(userID uuid.UUID) (dto.ResponseGetUserEmployment, error)
	GetUserSeniority(userID uuid.UUID) (dto.ResponseGetUserSeniority, error)
	GetUserWorkExperience(userID uuid.UUID) (dto.ResponseGetUserWorkExperience, error)
	GetUserHardSkill(userID uuid.UUID) ([]dto.ResponseGetUserHardSkill, error)
	GetUserSoftSkill(userID uuid.UUID) ([]dto.ResponseGetUserSoftSkill, error)
	GetUserTools(userID uuid.UUID) ([]dto.ResponseGetUserTools, error)
	GetUserLink(userID uuid.UUID) (dto.ResponseGetUserLink, error)
	GetUserSubscription(userID uuid.UUID) (dto.ResponseGetUserSubscription, error)
	UpdateUserInfo(updateUserInfo dto.UpdateUserInfo) (dto.ResponseUpdateUserInfo, error)
	UpdateUserDetail(updateUserDetail dto.UpdateUserDetail) (dto.ResponseUpdateUserDetail, error)
	UpdateUserContact(updateUserContact dto.UpdateUserContact) (dto.ResponseUpdateUserContact, error)
	UpdateUserEducation(updateUserEducation dto.UpdateUserEducation) (dto.ResponseUpdateUserEducation, error)
	AddUserLanguage(addUserLanguage dto.AddUserLanguage) ([]dto.ResponseGetUserLanguage, error)
	UpdateUserEmployment(updateUserEmployment dto.UpdateUserEmployment) (dto.ResponseUpdateUserEmployment, error)
	UpdateUserSeniority(updateUserSeniority dto.UpdateUserSeniority) (dto.ResponseUpdateUserSeniority, error)
	UpdateUserWorkExperience(updateUserWorkExperience dto.UpdateUserWorkExperience) (dto.ResponseUpdateUserWorkExperience, error)
	AddUserHardSkill(addUserHardSkill dto.AddUserHardSkill) ([]dto.ResponseGetUserHardSkill, error)
	AddUserSoftSkill(addUserSoftSkill dto.AddUserSoftSkill) ([]dto.ResponseGetUserSoftSkill, error)
	AddUserTools(addUserTools dto.AddUserTools) ([]dto.ResponseGetUserTools, error)
	UpdateUserLink(updateUserLink dto.UpdateUserLink) (dto.ResponseUpdateUserLink, error)
	DeleteUserLanguage(deleteUserLanguage dto.DeleteUserLanguage) error
	DeleteUserHardSkill(deleteUserHardSkill dto.DeleteUserHardSkill) error
	DeleteUserSoftSkill(deleteUserSoftSkill dto.DeleteUserSoftSkill) error
	DeleteUserTools(deleteUserTools dto.DeleteUserTools) error
	SoftDelete(userID uuid.UUID) error
}

type UserUseCase struct {
	userRepo     repository.UserDBItf
	jwt          jwt.JWTItf
	redis        redisitf.RedisItf
	redisContext context.Context
	s3           s3itf.S3Itf
	s3Context    context.Context
	env          *env.Env
}

func NewUserUseCase(
	userRepo repository.UserDBItf, jwt *jwt.JWT,
	redis redisitf.RedisItf, s3 s3itf.S3Itf,
	env *env.Env,
) UserUseCaseItf {
	return &UserUseCase{
		userRepo:     userRepo,
		jwt:          jwt,
		redis:        redis,
		redisContext: context.Background(),
		s3:           s3,
		s3Context:    context.Background(),
		env:          env,
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

	err = u.userRepo.Register(&user)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	return user.ParseToDTOResponseRegister(), nil
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

func (u *UserUseCase) UploadProfilePicture(userID uuid.UUID, file multipart.FileHeader) (dto.ResponseGetUserInfo, error) {
	fileContent, err := file.Open()
	if err != nil {
		return dto.ResponseGetUserInfo{}, fiber.ErrBadRequest
	}

	byteContainer, err := io.ReadAll(fileContent)
	if err != nil {
		return dto.ResponseGetUserInfo{}, fiber.ErrBadRequest
	}

	objectKey := fmt.Sprintf(
		"%s/%s",
		constants.ProfileDirectory,
		userID.String(),
	)

	profilePictureURL := fmt.Sprintf(
		"%s/%s",
		u.env.S3URL,
		objectKey,
	)

	err = u.s3.Upload(u.s3Context, objectKey, byteContainer)
	if err != nil {
		log.Println(err)
	}

	err = u.userRepo.UploadProfilePicture(userID, profilePictureURL)

	redisKey := fmt.Sprintf("user:%s", userID.String())
	u.redis.Delete(redisKey)

	return u.GetUserInfo(userID)
}

func (u *UserUseCase) GetUserIDFromUsername(username string) (uuid.UUID, error) {
	// TODO: remove
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
		return dto.ResponseGetUserInfo{}, err
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

func (u *UserUseCase) GetUserDetail(userID uuid.UUID) (dto.ResponseGetUserDetail, error) {
	userDetail := entity.UserDetail{
		UserID: userID,
	}

	err := u.userRepo.GetUserDetail(&userDetail)
	if err != nil {
		return dto.ResponseGetUserDetail{}, err
	}

	return userDetail.ParseToDTOResponseGetUserDetail(), nil
}

func (u *UserUseCase) GetUserContact(userID uuid.UUID) (dto.ResponseGetUserContact, error) {
	userContact := entity.UserContact{
		UserID: userID,
	}

	err := u.userRepo.GetUserContact(&userContact)
	if err != nil {
		return dto.ResponseGetUserContact{}, err
	}

	return userContact.ParseToDTOResponseGetUserContact(), nil
}

func (u *UserUseCase) GetUserEducation(userID uuid.UUID) (dto.ResponseGetUserEducation, error) {
	userEducation := entity.UserEducation{
		UserID: userID,
	}

	err := u.userRepo.GetUserEducation(&userEducation)
	if err != nil {
		return dto.ResponseGetUserEducation{}, err
	}

	return userEducation.ParseToDTOResponseGetUserEducation(), nil
}

func (u *UserUseCase) GetUserLanguage(userID uuid.UUID) ([]dto.ResponseGetUserLanguage, error) {
	var userLanguage []entity.UserLanguage

	err := u.userRepo.GetUserLanguage(&userLanguage)
	if err != nil {
		return nil, err
	}

	userLanguageList := make([]dto.ResponseGetUserLanguage, len(userLanguage))

	for i, userLanguageItem := range userLanguage {
		userLanguageList[i] = userLanguageItem.ParseToDTOResponseGetUserLanguage()
	}

	return userLanguageList, nil
}

func (u *UserUseCase) GetUserEmployment(userID uuid.UUID) (dto.ResponseGetUserEmployment, error) {
	userEmployment := entity.UserEmployment{
		UserID: userID,
	}

	err := u.userRepo.GetUserEmployment(&userEmployment)
	if err != nil {
		return dto.ResponseGetUserEmployment{}, err
	}

	return userEmployment.ParseToDTOResponseGetUserEmployment(), nil
}

func (u *UserUseCase) GetUserSeniority(userID uuid.UUID) (dto.ResponseGetUserSeniority, error) {
	userSeniority := entity.UserSeniority{
		UserID: userID,
	}

	err := u.userRepo.GetUserSeniority(&userSeniority)
	if err != nil {
		return dto.ResponseGetUserSeniority{}, err
	}

	return userSeniority.ParseToDTOResponseGetUserSeniority(), nil
}

func (u *UserUseCase) GetUserWorkExperience(userID uuid.UUID) (dto.ResponseGetUserWorkExperience, error) {
	userWorkExperience := entity.UserWorkExperience{
		UserID: userID,
	}

	err := u.userRepo.GetUserWorkExperience(&userWorkExperience)
	if err != nil {
		return dto.ResponseGetUserWorkExperience{}, err
	}

	return userWorkExperience.ParseToDTOResponseGetUserWorkExperience(), nil
}

func (u *UserUseCase) GetUserHardSkill(userID uuid.UUID) ([]dto.ResponseGetUserHardSkill, error) {
	var userHardSkill []entity.UserHardSkill

	err := u.userRepo.GetUserHardSkill(&userHardSkill)
	if err != nil {
		return nil, err
	}

	userHardSkillList := make([]dto.ResponseGetUserHardSkill, len(userHardSkill))

	for i, userHardSkillItem := range userHardSkill {
		userHardSkillList[i] = userHardSkillItem.ParseToDTOResponseGetUserHardSkill()
	}

	return userHardSkillList, nil
}

func (u *UserUseCase) GetUserSoftSkill(userID uuid.UUID) ([]dto.ResponseGetUserSoftSkill, error) {
	var userSoftSkill []entity.UserSoftSkill

	err := u.userRepo.GetUserSoftSkill(&userSoftSkill)
	if err != nil {
		return nil, err
	}

	userSoftSkillList := make([]dto.ResponseGetUserSoftSkill, len(userSoftSkill))

	for i, userSoftSkillItem := range userSoftSkill {
		userSoftSkillList[i] = userSoftSkillItem.ParseToDTOResponseGetUserSoftSkill()
	}

	return userSoftSkillList, nil
}

func (u *UserUseCase) GetUserTools(userID uuid.UUID) ([]dto.ResponseGetUserTools, error) {
	var userTools []entity.UserTools

	err := u.userRepo.GetUserTools(&userTools)
	if err != nil {
		return nil, err
	}

	userToolsList := make([]dto.ResponseGetUserTools, len(userTools))

	for i, userToolsItem := range userTools {
		userToolsList[i] = userToolsItem.ParseToDTOResponseGetUserTools()
	}

	return userToolsList, nil
}

func (u *UserUseCase) GetUserLink(userID uuid.UUID) (dto.ResponseGetUserLink, error) {
	userLink := entity.UserLink{
		UserID: userID,
	}

	err := u.userRepo.GetUserLink(&userLink)
	if err != nil {
		return dto.ResponseGetUserLink{}, err
	}

	return userLink.ParseToDTOResponseGetUserLink(), nil
}

func (u *UserUseCase) GetUserSubscription(userID uuid.UUID) (dto.ResponseGetUserSubscription, error) {
	userSubscription := entity.UserSubscription{
		UserID: userID,
	}

	err := u.userRepo.GetUserSubscription(&userSubscription)
	if err != nil {
		return dto.ResponseGetUserSubscription{}, err
	}

	return userSubscription.ParseToDTOResponseGetUserSubscription(), nil
}

func (u *UserUseCase) UpdateUserInfo(updateUserInfo dto.UpdateUserInfo) (dto.ResponseUpdateUserInfo, error) {
	var redisKey string

	user := entity.User{
		ID:       updateUserInfo.ID,
		Email:    updateUserInfo.Email,
		Username: updateUserInfo.Username,
		Name:     updateUserInfo.Name,
	}

	err := u.userRepo.UpdateUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	err = u.userRepo.GetUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserInfo.ID.String())

	go func() {
		newData, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		u.redis.Set(redisKey, string(newData))
	}()

	return user.ParseToDTOResponseUpdateUserInfo(), nil
}

func (u *UserUseCase) UpdateUserDetail(updateUserDetail dto.UpdateUserDetail) (dto.ResponseUpdateUserDetail, error) {
	var redisKey string

	userDetail := entity.UserDetail{
		UserID:      updateUserDetail.UserID,
		FirstName:   updateUserDetail.FirstName,
		LastName:    updateUserDetail.LastName,
		NickName:    updateUserDetail.NickName,
		Gender:      updateUserDetail.Gender,
		DateOfBirth: updateUserDetail.DateOfBirth,
		Nationality: updateUserDetail.Nationality,
		Location:    updateUserDetail.Location,
	}

	err := u.userRepo.UpdateUserDetail(&userDetail)
	if err != nil {
		return dto.ResponseUpdateUserDetail{},
			err
	}

	err = u.userRepo.GetUserDetail(&userDetail)
	if err != nil {
		return dto.ResponseUpdateUserDetail{},
			err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserDetail.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userDetail.ParseToDTOResponseUpdateUserDetail(), nil
}

func (u *UserUseCase) UpdateUserContact(updateUserContact dto.UpdateUserContact) (dto.ResponseUpdateUserContact, error) {
	var redisKey string

	userContact := entity.UserContact{
		UserID:           updateUserContact.UserID,
		AlternativeEmail: updateUserContact.AlternativeEmail,
		PhoneNumber:      updateUserContact.PhoneNumber,
		WhatsappNumber:   updateUserContact.WhatsappNumber,
	}

	err := u.userRepo.UpdateUserContact(&userContact)
	if err != nil {
		return dto.ResponseUpdateUserContact{}, err
	}

	err = u.userRepo.GetUserContact(&userContact)
	if err != nil {
		return dto.ResponseUpdateUserContact{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserContact.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userContact.ParseToDTOResponseUpdateUserContact(), nil
}

func (u *UserUseCase) UpdateUserEducation(updateUserEducation dto.UpdateUserEducation) (dto.ResponseUpdateUserEducation, error) {
	var redisKey string

	userEducation := entity.UserEducation{
		UserID:        updateUserEducation.UserID,
		LastEducation: updateUserEducation.LastEducation,
		Status:        updateUserEducation.Status,
		YearStarted:   updateUserEducation.YearStarted,
		YearEnded:     updateUserEducation.YearEnded,
	}

	err := u.userRepo.UpdateUserEducation(&userEducation)
	if err != nil {
		return dto.ResponseUpdateUserEducation{}, err
	}

	err = u.userRepo.GetUserEducation(&userEducation)
	if err != nil {
		return dto.ResponseUpdateUserEducation{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserEducation.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userEducation.ParseToDTOResponseUpdateUserEducation(), nil
}

func (u *UserUseCase) AddUserLanguage(addUserLanguage dto.AddUserLanguage) ([]dto.ResponseGetUserLanguage, error) {
	userLanguage := entity.UserLanguage{
		UserID:         addUserLanguage.UserID,
		LanguageSpoken: addUserLanguage.LanguageSpoken,
	}

	err := u.userRepo.AddUserLanguage(&userLanguage)
	if err != nil {
		return nil, err
	}

	return u.GetUserLanguage(addUserLanguage.UserID)
}

func (u *UserUseCase) UpdateUserEmployment(updateUserEmployment dto.UpdateUserEmployment) (dto.ResponseUpdateUserEmployment, error) {
	var redisKey string

	userEmployment := entity.UserEmployment{
		UserID:              updateUserEmployment.UserID,
		CurrentStatus:       updateUserEmployment.CurrentStatus,
		TotalWorkExperience: updateUserEmployment.TotalWorkExperience,
	}

	err := u.userRepo.UpdateUserEmployment(&userEmployment)
	if err != nil {
		return dto.ResponseUpdateUserEmployment{}, err
	}

	err = u.userRepo.GetUserEmployment(&userEmployment)
	if err != nil {
		return dto.ResponseUpdateUserEmployment{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserEmployment.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userEmployment.ParseToDTOResponseUpdateUserEmployment(), nil
}

func (u *UserUseCase) UpdateUserSeniority(updateUserSeniority dto.UpdateUserSeniority) (dto.ResponseUpdateUserSeniority, error) {
	var redisKey string

	userSeniority := entity.UserSeniority{
		UserID: updateUserSeniority.UserID,
		Year:   updateUserSeniority.Year,
	}

	err := u.userRepo.UpdateUserSeniority(&userSeniority)
	if err != nil {
		return dto.ResponseUpdateUserSeniority{}, err
	}

	err = u.userRepo.GetUserSeniority(&userSeniority)
	if err != nil {
		return dto.ResponseUpdateUserSeniority{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserSeniority.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userSeniority.ParseToDTOResponseUpdateUserSeniority(), nil
}

func (u *UserUseCase) UpdateUserWorkExperience(updateUserWorkExperience dto.UpdateUserWorkExperience) (dto.ResponseUpdateUserWorkExperience, error) {
	var redisKey string

	userWorkExperience := entity.UserWorkExperience{
		UserID:         updateUserWorkExperience.UserID,
		JobTitle:       updateUserWorkExperience.JobTitle,
		CompanyName:    updateUserWorkExperience.CompanyName,
		Industry:       updateUserWorkExperience.Industry,
		EmploymentType: updateUserWorkExperience.EmploymentType,
		StartDate:      updateUserWorkExperience.StartDate,
		EndDate:        updateUserWorkExperience.EndDate,
	}

	err := u.userRepo.UpdateUserWorkExperience(&userWorkExperience)
	if err != nil {
		return dto.ResponseUpdateUserWorkExperience{}, err
	}

	err = u.userRepo.GetUserWorkExperience(&userWorkExperience)
	if err != nil {
		return dto.ResponseUpdateUserWorkExperience{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserWorkExperience.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userWorkExperience.ParseToDTOResponseUpdateUserWorkExperience(), nil
}

func (u *UserUseCase) AddUserHardSkill(addUserHardSkill dto.AddUserHardSkill) ([]dto.ResponseGetUserHardSkill, error) {
	userHardSkill := entity.UserHardSkill{
		UserID:    addUserHardSkill.UserID,
		HardSkill: addUserHardSkill.HardSkill,
	}

	err := u.userRepo.AddUserHardSkill(&userHardSkill)
	if err != nil {
		return nil, err
	}

	return u.GetUserHardSkill(addUserHardSkill.UserID)
}

func (u *UserUseCase) AddUserSoftSkill(addUserSoftSkill dto.AddUserSoftSkill) ([]dto.ResponseGetUserSoftSkill, error) {
	userSoftSkill := entity.UserSoftSkill{
		UserID:    addUserSoftSkill.UserID,
		SoftSkill: addUserSoftSkill.SoftSkill,
	}

	err := u.userRepo.AddUserSoftSkill(&userSoftSkill)
	if err != nil {
		return nil, err
	}

	return u.GetUserSoftSkill(addUserSoftSkill.UserID)
}

func (u *UserUseCase) AddUserTools(addUserTools dto.AddUserTools) ([]dto.ResponseGetUserTools, error) {
	userTools := entity.UserTools{
		UserID: addUserTools.UserID,
		Tools:  addUserTools.Tools,
	}

	err := u.userRepo.AddUserTools(&userTools)
	if err != nil {
		return nil, err
	}

	return u.GetUserTools(addUserTools.UserID)
}

func (u *UserUseCase) UpdateUserLink(updateUserLink dto.UpdateUserLink) (dto.ResponseUpdateUserLink, error) {
	var redisKey string

	userLink := entity.UserLink{
		UserID:    updateUserLink.UserID,
		LinkedIn:  updateUserLink.LinkedIn,
		Portfolio: updateUserLink.Portfolio,
		GitHub:    updateUserLink.GitHub,
		Other:     updateUserLink.Other,
	}

	err := u.userRepo.UpdateUserLink(&userLink)
	if err != nil {
		return dto.ResponseUpdateUserLink{}, err
	}

	err = u.userRepo.GetUserLink(&userLink)
	if err != nil {
		return dto.ResponseUpdateUserLink{}, err
	}

	redisKey = fmt.Sprintf("user:%s", updateUserLink.UserID.String())

	go func() {
		u.redis.Delete(redisKey)
	}()

	return userLink.ParseToDTOResponseUpdateUserLink(), nil
}

func (u *UserUseCase) DeleteUserLanguage(deleteUserLanguage dto.DeleteUserLanguage) error {
	userLanguage := entity.UserLanguage{
		UserID:         deleteUserLanguage.UserID,
		LanguageSpoken: deleteUserLanguage.LanguageSpoken,
	}

	err := u.userRepo.DeleteUserLanguage(&userLanguage)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) DeleteUserHardSkill(deleteUserHardSkill dto.DeleteUserHardSkill) error {
	userHardSkill := entity.UserHardSkill{
		UserID:    deleteUserHardSkill.UserID,
		HardSkill: deleteUserHardSkill.HardSkill,
	}

	return u.userRepo.DeleteUserHardSkill(&userHardSkill)
}

func (u *UserUseCase) DeleteUserSoftSkill(deleteUserSoftSkill dto.DeleteUserSoftSkill) error {
	userSoftSkill := entity.UserSoftSkill{
		UserID:    deleteUserSoftSkill.UserID,
		SoftSkill: deleteUserSoftSkill.SoftSkill,
	}

	return u.userRepo.DeleteUserSoftSkill(&userSoftSkill)
}

func (u *UserUseCase) DeleteUserTools(deleteUserTools dto.DeleteUserTools) error {
	userTools := entity.UserTools{
		UserID: deleteUserTools.UserID,
		Tools:  deleteUserTools.Tools,
	}

	return u.userRepo.DeleteUserTools(&userTools)
}

func (u *UserUseCase) SoftDelete(userID uuid.UUID) error {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.SoftDelete(&user)

	return err
}
