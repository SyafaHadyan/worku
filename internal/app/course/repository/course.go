// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
)

type CourseDBItf interface {
	GetCourseList(offset *int, limit *int, course *[]entity.Course) error
	GetCourseInfo(count *int64, course *entity.Course) error
	SearchCourse(query *string, course *[]entity.Course) error
}

type CourseDB struct {
	db *gorm.DB
}

func NewCourseDB(db *gorm.DB) CourseDBItf {
	return &CourseDB{
		db: db,
	}
}

func (r *CourseDB) GetCourseList(offset *int, limit *int, course *[]entity.Course) error {
	return r.db.Debug().
		Model(&course).
		Select("courses.id, courses.name, courses.cover_image, courses.price").
		Limit(*limit).
		Offset(*offset).
		Find(course).
		Error
}

func (r *CourseDB) GetCourseInfo(count *int64, course *entity.Course) error {
	return r.db.Debug().
		Model(&course).
		Select("courses.id, courses.name, courses.description, courses.cover_image, courses.price, courses.created_at, courses.updated_at").
		First(&course).
		Count(count).
		Error
}

func (r *CourseDB) SearchCourse(query *string, course *[]entity.Course) error {
	// TODO: limit search result
	return r.db.Debug().
		Model(&course).
		Raw(`
		SELECT courses.id, courses.name, courses.cover_image, courses.price
		FROM courses
		WHERE MATCH(name,description)
		AGAINST (? IN NATURAL LANGUAGE MODE)
		`,
			query).
		Scan(course).
		Error
}
