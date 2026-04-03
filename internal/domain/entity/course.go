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

type CourseCategory struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"type:nvarchar(256)"`
	Description string    `gorm:"type:nvarchar(2048)"`
	CoverImage  string    `gorm:"type:nvarchar(512)"`
}

type Course struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	CategoryID   uuid.UUID `gorm:"type:char(36)"`
	Name         string    `gorm:"type:nvarchar(256);index:idx_fulltext_search,class:FULLTEXT"`
	Description  string    `gorm:"type:nvarchar(2048);index:idx_fulltext_search,class:FULLTEXT"`
	CoverImage   string    `gorm:"type:nvarchar(512)"`
	CourseVideo  []CourseVideo
	CourseModule CourseModule
	CreatedAt    time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type CourseVideo struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	CourseID uuid.UUID `gorm:"type:char(36)"`
	VideoURL string    `gorm:"type:nvarchar(512)"`
}

type CourseModule struct {
	CourseID         uuid.UUID          `gorm:"type:char(36);primaryKey"`
	Name             string             `gorm:"type:nvarchar(128)"`
	Description      string             `gorm:"type:mediumtext"`
	CourseModuleItem []CourseModuleItem `gorm:"foreignKey:ModuleID;references:CourseID"`
}

type CourseModuleItem struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	ModuleID    uuid.UUID `gorm:"type:char(36)"`
	Name        string    `gorm:"type:nvarchar(128)"`
	Description string    `gorm:"type:mediumtext"`
}

func (c *CourseCategory) ParseToDTOResponseGetCourseCategory() dto.ResponseGetCourseCategory {
	var response dto.ResponseGetCourseCategory

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (c *Course) ParseToDTOResponseGetCourseList() dto.ResponseGetCourseList {
	var response dto.ResponseGetCourseList

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (c *Course) ParseToDTOResponseGetCourseInfo() dto.ResponseGetCourseInfo {
	var response dto.ResponseGetCourseInfo

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (c *Course) ParseToDTOResponseSearchCourse() dto.ResponseSearchCourse {
	var response dto.ResponseSearchCourse

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (c *CourseVideo) ParseToDTOResponseGetCourseVideo() dto.ResponseGetCourseVideo {
	var response dto.ResponseGetCourseVideo

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (c *CourseModule) ParseToDTOResponseGetCourseModule() dto.ResponseGetCourseModule {
	var response dto.ResponseGetCourseModule

	err := copier.Copy(&response, c)
	if err != nil {
		log.Println(err)
	}

	return response
}
