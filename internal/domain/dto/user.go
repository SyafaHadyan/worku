// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"required,email"`
	Username string    `json:"username" validate:"required,min=2,max=32"`
	Password string    `json:"password" validate:"required,min=4,max=512"`
	Name     string    `json:"name" validate:"omitempty,min=3,max=64"`
}

type ResponseRegister struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Username string `json:"username" validate:"omitempty,min=2,max=32"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=4,max=512"`
}

type ResponseLogin struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ResponseGoogleOAuth struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type ResponseLinkedInOAuth struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type UserDetail struct {
	// TODO: remove?
	UserID uuid.UUID `json:"user_id"`
}

type CheckUserID struct {
	// TODO: remove?
	ID uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
}

type GetUserID struct {
	// TODO: remove?
	Username string `json:"username" validate:"required,min=3,max=32"`
}

type ResponseGetUserInfo struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	ProfilePicture string    `json:"profile_picture"`
	UserDetail     struct {
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		NickName    string    `json:"nick_name"`
		Gender      string    `json:"gender"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Nationality string    `json:"nationality"`
		Location    string    `json:"location"`
	} `json:"user_detail"`
	UserContact struct {
		AlternativeEmail string `json:"alternative_email"`
		PhoneNumber      string `json:"phone_number"`
		WhatsappNumber   string `json:"whatsapp_number"`
	} `json:"user_contact"`
	UserEducation struct {
		LastEducation string `json:"last_education"`
		Status        string `json:"status"`
		YearStarted   uint   `json:"year_started"`
		YearEnded     uint   `json:"year_ended"`
	} `json:"user_education"`
	UserLanguage []struct {
		LanguageSpoken string `json:"language_spoken"`
	} `json:"user_language"`
	UserEmployment struct {
		CurrentStatus       string `json:"current_status"`
		TotalWorkExperience uint   `json:"total_work_experience"`
	} `json:"user_employment"`
	UserSeniority struct {
		Year uint `json:"year"`
	} `json:"user_seniority"`
	UserWorkExperience struct {
		JobTitle       string    `json:"job_title"`
		CompanyName    string    `json:"company_name"`
		Industry       string    `json:"industry"`
		EmploymentType string    `json:"employment_type"`
		StartDate      time.Time `json:"start_date"`
		EndDate        time.Time `json:"end_date"`
	} `json:"user_work_experience"`
	UserHardSkill []struct {
		HardSkill string `json:"hard_skill"`
	} `json:"user_hard_skill"`
	UserSoftSkill []struct {
		SoftSkill string `json:"soft_skill"`
	} `json:"user_soft_skill"`
	UserTools []struct {
		Tools string `json:"tools"`
	} `json:"user_tools"`
	UserLink struct {
		LinkedIn  string `json:"linkedin"`
		Portfolio string `json:"portfolio"`
		GitHub    string `json:"github"`
		Other     string `json:"other"`
	} `json:"user_link"`
	UserSubscription struct {
		ExpiryDate time.Time `json:"expiry_date"`
	} `json:"user_subscription"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetUserDetail struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	NickName    string    `json:"nick_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Nationality string    `json:"nationality"`
	Location    string    `json:"location"`
}

type ResponseGetUserContact struct {
	AlternativeEmail string `json:"alternative_email"`
	PhoneNumber      string `json:"phone_number"`
	WhatsappNumber   string `json:"whatsapp_number"`
}

type ResponseGetUserEducation struct {
	LastEducation string `json:"last_education"`
	Status        string `json:"status"`
	YearStarted   uint   `json:"year_started"`
	YearEnded     uint   `json:"year_ended"`
}

type ResponseGetUserLanguage struct {
	LanguageSpoken string `json:"language_spoken"`
}

type ResponseGetUserEmployment struct {
	CurrentStatus       string `json:"current_status"`
	TotalWorkExperience uint   `json:"total_work_experience"`
}

type ResponseGetUserSeniority struct {
	Year uint `json:"year"`
}

type ResponseGetUserWorkExperience struct {
	JobTitle       string    `json:"job_title"`
	CompanyName    string    `json:"company_name"`
	Industry       string    `json:"industry"`
	EmploymentType string    `json:"employment_type"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

type ResponseGetUserHardSkill struct {
	HardSkill string `json:"hard_skill"`
}

type ResponseGetUserSoftSkill struct {
	SoftSkill string `json:"soft_skill"`
}

type ResponseGetUserTools struct {
	Tools string `json:"tools"`
}

type ResponseGetUserLink struct {
	LinkedIn  string `json:"linkedin"`
	Portfolio string `json:"portfolio"`
	GitHub    string `json:"github"`
	Other     string `json:"other"`
}

type ResponseGetUserSubscription struct {
	ExpiryDate time.Time `json:"expiry_date"`
}

type UpdateUserInfo struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"omitempty,email"`
	Username string    `json:"username" validate:"omitempty,min=3,max=32"`
	Password string    `json:"password" validate:"omitempty,min=4"`
	Name     string    `json:"name" validate:"omitempty,min=3,max=128"`
}

type ResponseUpdateUserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserDetail struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name" validate:"omitempty,min=2,max=128"`
	LastName    string    `json:"last_name" validate:"omitempty,min=2,max=128"`
	NickName    string    `json:"nick_name" validate:"omitempty,min=2,max=128"`
	Gender      string    `json:"gender" validate:"omitempty,min=1,max=1"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"omitempty"`
	Nationality string    `json:"nationality" validate:"omitempty,min=2,max=128"`
	Location    string    `json:"location" validate:"omitempty,min=2,max=128"`
}

type ResponseUpdateUserDetail struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	NickName    string    `json:"nick_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Nationality string    `json:"nationality"`
	Location    string    `json:"location"`
}

type UpdateUserContact struct {
	UserID           uuid.UUID `json:"user_id"`
	AlternativeEmail string    `json:"alternative_email" validate:"omitempty,email,max=256"`
	PhoneNumber      string    `json:"phone_number" validate:"omitempty,min=6,max=64"`
	WhatsappNumber   string    `json:"whatsapp_number" validate:"omitempty,min=6,max=64"`
}

type ResponseUpdateUserContact struct {
	AlternativeEmail string `json:"alternative_email"`
	PhoneNumber      string `json:"phone_number"`
	WhatsappNumber   string `json:"whatsapp_number"`
}

type UpdateUserEducation struct {
	UserID        uuid.UUID `json:"user_id"`
	LastEducation string    `json:"last_education" validate:"omitempty,min=2,max=128"`
	Status        string    `json:"status" validate:"omitempty,min=2,max=128"`
	YearStarted   uint      `json:"year_started" validate:"omitempty,numeric"`
	YearEnded     uint      `json:"year_ended" validate:"omitempty,numeric"`
}

type ResponseUpdateUserEducation struct {
	LastEducation string `json:"last_education"`
	Status        string `json:"status"`
	YearStarted   uint   `json:"year_started"`
	YearEnded     uint   `json:"year_ended"`
}

type AddUserLanguage struct {
	UserID         uuid.UUID `json:"user_id"`
	LanguageSpoken string    `json:"language_spoken" validate:"required,min=2,max=128"`
}

type DeleteUserLanguage struct {
	UserID         uuid.UUID `json:"user_id"`
	LanguageSpoken string    `json:"language_spoken" validate:"required,min=2,max=128"`
}

type UpdateUserEmployment struct {
	UserID              uuid.UUID `json:"user_id"`
	CurrentStatus       string    `json:"current_status" validate:"omitempty,min=2,max=182"`
	TotalWorkExperience uint      `json:"total_work_experience" validate:"omitempty,numeric"`
}

type ResponseUpdateUserEmployment struct {
	CurrentStatus       string `json:"current_status"`
	TotalWorkExperience uint   `json:"total_work_experience"`
}

type UpdateUserSeniority struct {
	UserID uuid.UUID `json:"user_id"`
	Year   uint      `json:"year" validate:"omitempty,numeric"`
}

type ResponseUpdateUserSeniority struct {
	Year uint `json:"year"`
}

type UpdateUserWorkExperience struct {
	UserID         uuid.UUID `json:"user_id"`
	JobTitle       string    `json:"job_title" validate:"omitempty,min=2,max=256"`
	CompanyName    string    `json:"company_name" validate:"omitempty,min=2,max=256"`
	Industry       string    `json:"industry" validate:"omitempty,min=2,max=256"`
	EmploymentType string    `json:"employment_type" validate:"omitempty,min=2,max=256"`
	StartDate      time.Time `json:"start_date" validate:"omitempty"`
	EndDate        time.Time `json:"end_date" validate:"omitempty"`
}

type ResponseUpdateUserWorkExperience struct {
	JobTitle       string    `json:"job_title"`
	CompanyName    string    `json:"company_name"`
	Industry       string    `json:"industry"`
	EmploymentType string    `json:"employment_type"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

type AddUserHardSkill struct {
	UserID    uuid.UUID `json:"user_id"`
	HardSkill string    `json:"hard_skill" validate:"required,min=2,max=256"`
}

type DeleteUserHardSkill struct {
	UserID    uuid.UUID `json:"user_id"`
	HardSkill string    `json:"hard_skill" validate:"required,min=2,max=256"`
}

type AddUserSoftSkill struct {
	UserID    uuid.UUID `json:"user_id"`
	SoftSkill string    `json:"soft_skill" validate:"required,min=2,max=256"`
}

type DeleteUserSoftSkill struct {
	UserID    uuid.UUID `json:"user_id"`
	SoftSkill string    `json:"soft_skill" validate:"required,min=2,max=256"`
}

type AddUserTools struct {
	UserID uuid.UUID `json:"user_id"`
	Tools  string    `json:"tools" validate:"required,min=2,max=256"`
}

type DeleteUserTools struct {
	UserID uuid.UUID `json:"user_id"`
	Tools  string    `json:"tools" validate:"required,min=2,max=256"`
}

type UpdateUserLink struct {
	UserID    uuid.UUID `json:"user_id"`
	LinkedIn  string    `json:"linkedin" validate:"omitempty,url,max=256"`
	Portfolio string    `json:"portfolio" validate:"omitempty,url,max=256"`
	GitHub    string    `json:"github" validate:"omitempty,url,max=256"`
	Other     string    `json:"other" validate:"omitempty,url,max=256"`
}

type ResponseUpdateUserLink struct {
	LinkedIn  string `json:"linkedin"`
	Portfolio string `json:"portfolio"`
	GitHub    string `json:"github"`
	Other     string `json:"other"`
}
