// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDBItf interface {
	Register(user *entity.User) error
	Login(user *entity.User) error
	GoogleOAuthCreateUser(user *entity.User) error
	GoogleOAuthCheckUser(user *entity.User) error
	LinkedInOAuthCreateUser(user *entity.User) error
	LinkedInOAuthCheckUser(user *entity.User) error
	UploadProfilePicture(userID uuid.UUID, profilePictureURL string) error
	UpdateUserInfo(user *entity.User) error
	UpdateUserDetail(userDetail *entity.UserDetail) error
	UpdateUserContact(userContact *entity.UserContact) error
	UpdateUserEducation(userEducation *entity.UserEducation) error
	AddUserLanguage(userLanguage *entity.UserLanguage) error
	UpdateUserEmployment(userEmployment *entity.UserEmployment) error
	UpdateUserSeniority(userSeniority *entity.UserSeniority) error
	UpdateUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error
	AddUserHardSkill(userHardSkill *entity.UserHardSkill) error
	AddUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error
	AddUserTools(userTools *entity.UserTools) error
	UpdateUserLink(userLink *entity.UserLink) error
	GetUserInfo(user *entity.User) error
	GetUserDetail(userDetail *entity.UserDetail) error
	GetUserContact(userContact *entity.UserContact) error
	GetUserEducation(userEducation *entity.UserEducation) error
	GetUserLanguage(userID uuid.UUID, userLanguage *[]entity.UserLanguage) error
	GetUserEmployment(userEmployment *entity.UserEmployment) error
	GetUserSeniority(userSeniority *entity.UserSeniority) error
	GetUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error
	GetUserHardSkill(userID uuid.UUID, userHardSkill *[]entity.UserHardSkill) error
	GetUserSoftSkill(userID uuid.UUID, userSoftSkill *[]entity.UserSoftSkill) error
	GetUserTools(userID uuid.UUID, userTools *[]entity.UserTools) error
	GetUserLink(userLink *entity.UserLink) error
	GetUserSubscription(userSubscription *entity.UserSubscription) error
	DeleteUserLanguage(userLanguage *entity.UserLanguage) error
	DeleteUserHardSkill(userHardSkill *entity.UserHardSkill) error
	DeleteUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error
	DeleteUserTools(userTools *entity.UserTools) error
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
	return r.db.
		Model(&entity.User{}).
		Create(user).
		Error
}

func (r *UserDB) Login(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Where("username = ? OR email = ?", user.Username, user.Email).
		First(user).
		Error
}

func (r *UserDB) GoogleOAuthCreateUser(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Create(user).
		Error
}

func (r *UserDB) GoogleOAuthCheckUser(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Where("email = ?", user.Email).
		First(user).
		Error
}

func (r *UserDB) LinkedInOAuthCreateUser(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Create(user).
		Error
}

func (r *UserDB) LinkedInOAuthCheckUser(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Where("email = ?", user.Email).
		First(user).
		Error
}

func (r *UserDB) UploadProfilePicture(userID uuid.UUID, profilePictureURL string) error {
	return r.db.
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("profile_picture", profilePictureURL).
		Error
}

func (r *UserDB) UpdateUserInfo(user *entity.User) error {
	return r.db.
		Model(&entity.User{}).
		Updates(user).
		Error
}

func (r *UserDB) UpdateUserDetail(userDetail *entity.UserDetail) error {
	var err error

	if r.db.
		Model(&entity.UserDetail{}).
		Where("user_id = ?", userDetail.UserID).
		Updates(userDetail).
		RowsAffected == 0 {
		err = r.db.
			Create(userDetail).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserContact(userContact *entity.UserContact) error {
	var err error

	if r.db.
		Model(&entity.UserContact{}).
		Where("user_id = ?", userContact.UserID).
		Updates(userContact).
		RowsAffected == 0 {
		err = r.db.
			Create(userContact).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserEducation(userEducation *entity.UserEducation) error {
	var err error

	if r.db.
		Model(&entity.UserEducation{}).
		Where("user_id = ?", userEducation.UserID).
		Updates(userEducation).
		RowsAffected == 0 {
		err = r.db.
			Create(userEducation).
			Error
	}

	return err
}

func (r *UserDB) AddUserLanguage(userLanguage *entity.UserLanguage) error {
	return r.db.
		Model(&entity.UserLanguage{}).
		Create(userLanguage).
		Error
}

func (r *UserDB) UpdateUserEmployment(userEmployment *entity.UserEmployment) error {
	var err error

	if r.db.
		Model(&entity.UserEmployment{}).
		Where("user_id = ?", userEmployment.UserID).
		Updates(userEmployment).
		RowsAffected == 0 {
		err = r.db.
			Create(userEmployment).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserSeniority(userSeniority *entity.UserSeniority) error {
	var err error

	if r.db.
		Model(&entity.UserSeniority{}).
		Where("user_id = ?", userSeniority.UserID).
		Updates(userSeniority).
		RowsAffected == 0 {
		err = r.db.
			Create(userSeniority).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error {
	var err error

	if r.db.
		Model(&entity.UserWorkExperience{}).
		Where("user_id = ?", userWorkExperience.UserID).
		Updates(userWorkExperience).
		RowsAffected == 0 {
		err = r.db.
			Create(userWorkExperience).
			Error
	}

	return err
}

func (r *UserDB) AddUserHardSkill(userHardSkill *entity.UserHardSkill) error {
	return r.db.
		Model(&entity.UserHardSkill{}).
		Create(userHardSkill).
		Error
}

func (r *UserDB) AddUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error {
	return r.db.
		Model(&entity.UserSoftSkill{}).
		Create(userSoftSkill).
		Error
}

func (r *UserDB) AddUserTools(userTools *entity.UserTools) error {
	return r.db.
		Model(&entity.UserTools{}).
		Create(userTools).
		Error
}

func (r *UserDB) UpdateUserLink(userLink *entity.UserLink) error {
	var err error

	if r.db.
		Model(&entity.UserLink{}).
		Where("user_id = ?", userLink.UserID).
		Updates(userLink).
		RowsAffected == 0 {
		err = r.db.
			Create(userLink).
			Error
	}

	return err
}

func (r *UserDB) GetUserInfo(user *entity.User) error {
	return r.db.
		Preload(clause.Associations).
		First(user).
		Error
}

func (r *UserDB) GetUserDetail(userDetail *entity.UserDetail) error {
	return r.db.
		Model(&entity.UserDetail{}).
		First(userDetail).
		Error
}

func (r *UserDB) GetUserContact(userContact *entity.UserContact) error {
	return r.db.
		Model(&entity.UserContact{}).
		First(userContact).
		Error
}

func (r *UserDB) GetUserEducation(userEducation *entity.UserEducation) error {
	return r.db.
		Model(&entity.UserEducation{}).
		First(userEducation).
		Error
}

func (r *UserDB) GetUserLanguage(userID uuid.UUID, userLanguage *[]entity.UserLanguage) error {
	return r.db.
		Model(&entity.UserLanguage{}).
		Where("user_id = ?", userID).
		Find(userLanguage).
		Error
}

func (r *UserDB) GetUserEmployment(userEmployment *entity.UserEmployment) error {
	return r.db.
		Model(&entity.UserEmployment{}).
		First(userEmployment).
		Error
}

func (r *UserDB) GetUserSeniority(userSeniority *entity.UserSeniority) error {
	return r.db.
		Model(&entity.UserSeniority{}).
		First(userSeniority).
		Error
}

func (r *UserDB) GetUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error {
	return r.db.
		Model(&entity.UserWorkExperience{}).
		First(userWorkExperience).
		Error
}

func (r *UserDB) GetUserHardSkill(userID uuid.UUID, userHardSkill *[]entity.UserHardSkill) error {
	return r.db.
		Model(&entity.UserHardSkill{}).
		Where("user_id = ?", userID).
		Find(userHardSkill).
		Error
}

func (r *UserDB) GetUserSoftSkill(userID uuid.UUID, userSoftSkill *[]entity.UserSoftSkill) error {
	return r.db.
		Model(&entity.UserSoftSkill{}).
		Where("user_id = ?", userID).
		Find(userSoftSkill).
		Error
}

func (r *UserDB) GetUserTools(userID uuid.UUID, userTools *[]entity.UserTools) error {
	return r.db.
		Model(&entity.UserTools{}).
		Where("user_id = ?", userID).
		Find(userTools).
		Error
}

func (r *UserDB) GetUserLink(userLink *entity.UserLink) error {
	return r.db.
		Model(&entity.UserLink{}).
		First(userLink).
		Error
}

func (r *UserDB) GetUserSubscription(userSubscription *entity.UserSubscription) error {
	return r.db.
		Model(&entity.UserSubscription{}).
		First(userSubscription).
		Error
}

func (r *UserDB) DeleteUserLanguage(userLanguage *entity.UserLanguage) error {
	rowsAffected := r.db.
		Model(&entity.UserLanguage{}).
		Delete(userLanguage).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserHardSkill(userHardSkill *entity.UserHardSkill) error {
	rowsAffected := r.db.
		Model(&entity.UserHardSkill{}).
		Delete(userHardSkill).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error {
	rowsAffected := r.db.
		Model(&entity.UserSoftSkill{}).
		Delete(userSoftSkill).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserTools(userTools *entity.UserTools) error {
	rowsAffected := r.db.
		Model(&entity.UserTools{}).
		Delete(userTools).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) SoftDelete(user *entity.User) error {
	return r.db.
		Delete(user).
		Error
}
