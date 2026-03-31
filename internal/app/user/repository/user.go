// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDBItf interface {
	Register(user *entity.User) error
	Login(user *entity.User) error
	GoogleOAuth(user *entity.User) error
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
	GetUserLanguage(userLanguage *[]entity.UserLanguage) error
	GetUserEmployment(userEmployment *entity.UserEmployment) error
	GetUserSeniority(userSeniority *entity.UserSeniority) error
	GetUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error
	GetUserHardSkill(userHardSkill *[]entity.UserHardSkill) error
	GetUserSoftSkill(userSoftSkill *[]entity.UserSoftSkill) error
	GetUserTools(userTools *[]entity.UserTools) error
	GetUserLink(userLink *entity.UserLink) error
	GetUserSubscription(userSubscription *entity.UserSubscription) error
	CheckUsername(user *entity.User) error
	GetUserIDFromUsername(user *entity.User) error
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
	return r.db.Debug().
		Model(&entity.User{}).
		Create(user).
		Error
}

func (r *UserDB) Login(user *entity.User) error {
	return r.db.Debug().
		Model(&entity.User{}).
		Where("username = ? OR email = ?", user.Username, user.Email).
		First(user).
		Error
}

func (r *UserDB) GoogleOAuth(user *entity.User) error {
	// TODO: Fix

	r.db.Debug().
		Model(&entity.User{}).
		Create(user)

	return r.db.Debug().
		Model(user).
		Where("users.email = ?", user.Email).
		First(user).
		Error
}

func (r *UserDB) UpdateUserInfo(user *entity.User) error {
	return r.db.Debug().
		Model(&entity.User{}).
		Updates(user).
		Error
}

func (r *UserDB) UpdateUserDetail(userDetail *entity.UserDetail) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserDetail{}).
		Where("user_id = ?", userDetail.UserID).
		Updates(userDetail).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userDetail).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserContact(userContact *entity.UserContact) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserContact{}).
		Where("user_id = ?", userContact.UserID).
		Updates(userContact).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userContact).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserEducation(userEducation *entity.UserEducation) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserEducation{}).
		Where("user_id = ?", userEducation.UserID).
		Updates(userEducation).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userEducation).
			Error
	}

	return err
}

func (r *UserDB) AddUserLanguage(userLanguage *entity.UserLanguage) error {
	return r.db.Debug().
		Model(&entity.UserLanguage{}).
		Create(userLanguage).
		Error
}

func (r *UserDB) UpdateUserEmployment(userEmployment *entity.UserEmployment) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserEmployment{}).
		Where("user_id = ?", userEmployment.UserID).
		Updates(userEmployment).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userEmployment).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserSeniority(userSeniority *entity.UserSeniority) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserSeniority{}).
		Where("user_id = ?", userSeniority.UserID).
		Updates(userSeniority).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userSeniority).
			Error
	}

	return err
}

func (r *UserDB) UpdateUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserWorkExperience{}).
		Where("user_id = ?", userWorkExperience.UserID).
		Updates(userWorkExperience).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userWorkExperience).
			Error
	}

	return err
}

func (r *UserDB) AddUserHardSkill(userHardSkill *entity.UserHardSkill) error {
	return r.db.Debug().
		Model(&entity.UserHardSkill{}).
		Create(userHardSkill).
		Error
}

func (r *UserDB) AddUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error {
	return r.db.Debug().
		Model(&entity.UserSoftSkill{}).
		Create(userSoftSkill).
		Error
}

func (r *UserDB) AddUserTools(userTools *entity.UserTools) error {
	return r.db.Debug().
		Model(&entity.UserTools{}).
		Create(userTools).
		Error
}

func (r *UserDB) UpdateUserLink(userLink *entity.UserLink) error {
	var err error

	if r.db.Debug().
		Model(&entity.UserLink{}).
		Where("user_id = ?", userLink.UserID).
		Updates(userLink).
		RowsAffected == 0 {
		err = r.db.Debug().
			Create(userLink).
			Error
	}

	return err
}

func (r *UserDB) CheckUsername(user *entity.User) error {
	// TODO: remove?
	return r.db.Debug().
		Raw("SELECT `username` FROM `users` WHERE username = ?", user.Username).
		First(&user).
		Error
}

func (r *UserDB) GetUserInfo(user *entity.User) error {
	return r.db.Debug().
		Preload(clause.Associations).
		First(user).
		Error
}

func (r *UserDB) GetUserDetail(userDetail *entity.UserDetail) error {
	return r.db.Debug().
		Model(&entity.UserDetail{}).
		First(userDetail).
		Error
}

func (r *UserDB) GetUserContact(userContact *entity.UserContact) error {
	return r.db.Debug().
		Model(&entity.UserContact{}).
		First(userContact).
		Error
}

func (r *UserDB) GetUserEducation(userEducation *entity.UserEducation) error {
	return r.db.Debug().
		Model(&entity.UserEducation{}).
		First(userEducation).
		Error
}

func (r *UserDB) GetUserLanguage(userLanguage *[]entity.UserLanguage) error {
	return r.db.Debug().
		Model(&entity.UserLanguage{}).
		Find(userLanguage).
		Error
}

func (r *UserDB) GetUserEmployment(userEmployment *entity.UserEmployment) error {
	return r.db.Debug().
		Model(&entity.UserEmployment{}).
		First(userEmployment).
		Error
}

func (r *UserDB) GetUserSeniority(userSeniority *entity.UserSeniority) error {
	return r.db.Debug().
		Model(&entity.UserSeniority{}).
		First(userSeniority).
		Error
}

func (r *UserDB) GetUserWorkExperience(userWorkExperience *entity.UserWorkExperience) error {
	return r.db.Debug().
		Model(&entity.UserWorkExperience{}).
		First(userWorkExperience).
		Error
}

func (r *UserDB) GetUserHardSkill(userHardSkill *[]entity.UserHardSkill) error {
	return r.db.Debug().
		Model(&entity.UserHardSkill{}).
		Find(userHardSkill).
		Error
}

func (r *UserDB) GetUserSoftSkill(userSoftSkill *[]entity.UserSoftSkill) error {
	return r.db.Debug().
		Model(&entity.UserSoftSkill{}).
		Find(userSoftSkill).
		Error
}

func (r *UserDB) GetUserTools(userTools *[]entity.UserTools) error {
	return r.db.Debug().
		Model(&entity.UserTools{}).
		Find(userTools).
		Error
}

func (r *UserDB) GetUserLink(userLink *entity.UserLink) error {
	return r.db.Debug().
		Model(&entity.UserLink{}).
		First(userLink).
		Error
}

func (r *UserDB) GetUserSubscription(userSubscription *entity.UserSubscription) error {
	return r.db.Debug().
		Model(&entity.UserSubscription{}).
		First(userSubscription).
		Error
}

func (r *UserDB) GetUserIDFromUsername(user *entity.User) error {
	// TODO: remove
	return r.db.Debug().
		Select("id").
		Where("username = ?", user.Username).
		First(user).
		Error
}

func (r *UserDB) DeleteUserLanguage(userLanguage *entity.UserLanguage) error {
	rowsAffected := r.db.Debug().
		Model(&entity.UserLanguage{}).
		Delete(userLanguage).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserHardSkill(userHardSkill *entity.UserHardSkill) error {
	rowsAffected := r.db.Debug().
		Model(&entity.UserHardSkill{}).
		Delete(userHardSkill).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserSoftSkill(userSoftSkill *entity.UserSoftSkill) error {
	rowsAffected := r.db.Debug().
		Model(&entity.UserSoftSkill{}).
		Delete(userSoftSkill).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) DeleteUserTools(userTools *entity.UserTools) error {
	rowsAffected := r.db.Debug().
		Model(&entity.UserTools{}).
		Delete(userTools).
		RowsAffected

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserDB) SoftDelete(user *entity.User) error {
	return r.db.Debug().
		Delete(user).
		Error
}
