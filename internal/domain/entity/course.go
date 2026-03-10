// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name" gorm:"type:nvarchar(256);index:idx_fulltext_search,class:FULLTEXT"`
	Description string         `json:"description" gorm:"type:nvarchar(2048);index:idx_fulltext_search,class:FULLTEXT"`
	CoverImage  string         `json:"cover_image" gorm:"type:nvarchar(512)"`
	Price       uint32         `json:"price" gorm:"type:integer unsigned"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (c *Course) ParseToDTOResponseGetCourseList() dto.ResponseGetCourseList {
	return dto.ResponseGetCourseList{
		ID:         c.ID,
		Name:       c.Name,
		CoverImage: c.CoverImage,
		Price:      c.Price,
	}
}

func (c *Course) ParseToDTOResponseGetCourseInfo() dto.ResponseGetCourseInfo {
	return dto.ResponseGetCourseInfo{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CoverImage:  c.CoverImage,
		Price:       c.Price,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (c *Course) ParseToDTOResponseSearchCourse() dto.ResponseSearchCourse {
	return dto.ResponseSearchCourse{
		ID:         c.ID,
		Name:       c.Name,
		CoverImage: c.CoverImage,
		Price:      c.Price,
	}
}
