// Package entity defines database table and its relations
package entity

import (
	"log"
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type User struct {
	// TODO: remove name after fe confirms
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey"`
	Email              string    `gorm:"type:nvarchar(256);not null;unique"`
	Username           string    `gorm:"type:nvarchar(32);not null;unique"`
	Password           string    `gorm:"type:text;not null"`
	Name               string    `gorm:"type:nvarchar(128)"`
	UserDetail         UserDetail
	UserContact        UserContact
	UserEducation      UserEducation
	UserLanguage       []UserLanguage
	UserEmployment     UserEmployment
	UserSeniority      UserSeniority
	UserWorkExperience UserWorkExperience
	UserHardSkill      []UserHardSkill
	UserSoftSkill      []UserSoftSkill
	UserTools          []UserTools
	UserLink           UserLink
	UserSubscription   UserSubscription
	CreatedAt          time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type UserDetail struct {
	UserID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	FirstName   string    `gorm:"type:nvarchar(128)"`
	LastName    string    `gorm:"type:nvarchar(128)"`
	NickName    string    `gorm:"type:nvarchar(128)"`
	Gender      string    `gorm:"type:char(1)"`
	DateOfBirth time.Time `gorm:"type:date"`
	Nationality string    `gorm:"type:nvarchar(128)"`
	Location    string    `gorm:"type:nvarchar(128)"`
}

type UserContact struct {
	UserID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	AlternativeEmail string    `gorm:"type:nvarchar(256);unique"`
	PhoneNumber      string    `gorm:"type:nvarchar(64)"`
	WhatsappNumber   string    `gorm:"type:nvarchar(64)"`
}

type UserEducation struct {
	UserID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	LastEducation string    `gorm:"type:nvarchar(128)"`
	Status        string    `gorm:"type:nvarchar(128)"`
	YearStarted   uint      `gorm:"type:year"`
	YearEnded     uint      `gorm:"type:year"`
}

type UserLanguage struct {
	UserID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	LanguageSpoken string    `gorm:"type:nvarchar(128);primaryKey"`
}

type UserEmployment struct {
	UserID              uuid.UUID `gorm:"type:char(36);primaryKey"`
	CurrentStatus       string    `gorm:"type:nvarchar(182)"`
	TotalWorkExperience uint      `gorm:"type:integer unsigned"`
}

type UserSeniority struct {
	UserID uuid.UUID `gorm:"type:char(36);primaryKey"`
	Year   uint      `gorm:"type:integer unsigned"`
}

type UserWorkExperience struct {
	UserID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	JobTitle       string    `gorm:"type:nvarchar(256)"`
	CompanyName    string    `gorm:"type:nvarchar(256)"`
	Industry       string    `gorm:"type:nvarchar(256)"`
	EmploymentType string    `gorm:"type:nvarchar(256)"`
	StartDate      time.Time `gorm:"type:date"`
	EndDate        time.Time `gorm:"type:date"`
}

type UserHardSkill struct {
	UserID    uuid.UUID `gorm:"type:char(36);primaryKey"`
	HardSkill string    `gorm:"type:nvarchar(256);primaryKey"`
}

type UserSoftSkill struct {
	UserID    uuid.UUID `gorm:"type:char(36);primaryKey"`
	SoftSkill string    `gorm:"type:nvarchar(256);primaryKey"`
}

type UserTools struct {
	UserID uuid.UUID `gorm:"type:char(36);primaryKey"`
	Tools  string    `gorm:"type:nvarchar(256);primaryKey"`
}

type UserLink struct {
	UserID    uuid.UUID `gorm:"type:char(36);primaryKey"`
	LinkedIn  string    `gorm:"type:nvarchar(256)"`
	Portfolio string    `gorm:"type:nvarchar(256)"`
	GitHub    string    `gorm:"type:nvarchar(256)"`
	Other     string    `gorm:"type:nvarchar(256)"`
}

type UserSubscription struct {
	UserID     uuid.UUID `gorm:"type:char(36);primaryKey"`
	ExpiryDate time.Time `gorm:"type:timestamp"`
}

type UserCourse struct {
	UserID   uuid.UUID `gorm:"type:char(36);primaryKey"`
	CourseID string    `gorm:"type:char(36);primaryKey"`
}

func (u *User) ParseToDTOResponseRegister() dto.ResponseRegister {
	// TODO: change to copier
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
	// TODO: change to copier
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
	var response dto.ResponseGetUserInfo

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserDetail) ParseToDTOResponseGetUserDetail() dto.ResponseGetUserDetail {
	var response dto.ResponseGetUserDetail

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserContact) ParseToDTOResponseGetUserContact() dto.ResponseGetUserContact {
	var response dto.ResponseGetUserContact

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserEducation) ParseToDTOResponseGetUserEducation() dto.ResponseGetUserEducation {
	var response dto.ResponseGetUserEducation

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserLanguage) ParseToDTOResponseGetUserLanguage() dto.ResponseGetUserLanguage {
	var response dto.ResponseGetUserLanguage

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserEmployment) ParseToDTOResponseGetUserEmployment() dto.ResponseGetUserEmployment {
	var response dto.ResponseGetUserEmployment

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserSeniority) ParseToDTOResponseGetUserSeniority() dto.ResponseGetUserSeniority {
	var response dto.ResponseGetUserSeniority

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserWorkExperience) ParseToDTOResponseGetUserWorkExperience() dto.ResponseGetUserWorkExperience {
	var response dto.ResponseGetUserWorkExperience

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserHardSkill) ParseToDTOResponseGetUserHardSkill() dto.ResponseGetUserHardSkill {
	var response dto.ResponseGetUserHardSkill

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserSoftSkill) ParseToDTOResponseGetUserSoftSkill() dto.ResponseGetUserSoftSkill {
	var response dto.ResponseGetUserSoftSkill

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserTools) ParseToDTOResponseGetUserTools() dto.ResponseGetUserTools {
	var response dto.ResponseGetUserTools

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserLink) ParseToDTOResponseGetUserLink() dto.ResponseGetUserLink {
	var response dto.ResponseGetUserLink

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserSubscription) ParseToDTOResponseGetUserSubscription() dto.ResponseGetUserSubscription {
	var response dto.ResponseGetUserSubscription

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *User) ParseToDTOResponseUpdateUserInfo() dto.ResponseUpdateUserInfo {
	var response dto.ResponseUpdateUserInfo

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserDetail) ParseToDTOResponseUpdateUserDetail() dto.ResponseUpdateUserDetail {
	var response dto.ResponseUpdateUserDetail

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserContact) ParseToDTOResponseUpdateUserContact() dto.ResponseUpdateUserContact {
	var response dto.ResponseUpdateUserContact

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserEducation) ParseToDTOResponseUpdateUserEducation() dto.ResponseUpdateUserEducation {
	var response dto.ResponseUpdateUserEducation

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserEmployment) ParseToDTOResponseUpdateUserEmployment() dto.ResponseUpdateUserEmployment {
	var response dto.ResponseUpdateUserEmployment

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserSeniority) ParseToDTOResponseUpdateUserSeniority() dto.ResponseUpdateUserSeniority {
	var response dto.ResponseUpdateUserSeniority

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserWorkExperience) ParseToDTOResponseUpdateUserWorkExperience() dto.ResponseUpdateUserWorkExperience {
	var response dto.ResponseUpdateUserWorkExperience

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (u *UserLink) ParseToDTOResponseUpdateUserLink() dto.ResponseUpdateUserLink {
	var response dto.ResponseUpdateUserLink

	err := copier.Copy(&response, u)
	if err != nil {
		log.Println(err)
	}

	return response
}
