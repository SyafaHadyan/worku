// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourseDBItf interface {
	GetCourseCategory(courseCategory *[]entity.CourseCategory) error
	GetCourseList(offset *int, limit *int, course *[]entity.Course) error
	GetCourseListByCategory(categoryID uuid.UUID, offset *int, limit *int, course *[]entity.Course) error
	SearchCourse(offset *int, limit *int, query *string, course *[]entity.Course) error
	GetCourseInfo(course *entity.Course) error
	GetCourseVideo(courseID uuid.UUID, courseVideo *[]entity.CourseVideo) error
	GetCourseModule(courseID uuid.UUID, courseModule *[]entity.CourseModule) error
	GetCourseEnrollmentCount(courseID uuid.UUID) (int64, error)
	UpdateCourseEnrollment(userCourse *entity.UserCourse) error
}

type CourseDB struct {
	db *gorm.DB
}

func NewCourseDB(db *gorm.DB) CourseDBItf {
	return &CourseDB{
		db: db,
	}
}

func (r *CourseDB) GetCourseCategory(courseCategory *[]entity.CourseCategory) error {
	return r.db.
		Model(&entity.CourseCategory{}).
		Find(courseCategory).
		Error
}

func (r *CourseDB) GetCourseList(offset *int, limit *int, course *[]entity.Course) error {
	return r.db.
		Model(&entity.Course{}).
		Limit(*limit).
		Offset(*offset).
		Find(course).
		Error
}

func (r *CourseDB) GetCourseListByCategory(categoryID uuid.UUID, offset *int, limit *int, course *[]entity.Course) error {
	return r.db.
		Model(&entity.Course{}).
		Where("category_id = ?", categoryID).
		Limit(*limit).
		Offset(*offset).
		Find(course).
		Error
}

func (r *CourseDB) SearchCourse(offset *int, limit *int, query *string, course *[]entity.Course) error {
	return r.db.
		Model(&entity.Course{}).
		Raw(`
		SELECT * 
		FROM courses
		WHERE MATCH(name,description)
		AGAINST (? IN NATURAL LANGUAGE MODE)
		LIMIT ? OFFSET ? 
		`,
			query, limit, offset).
		Scan(course).
		Error
}

func (r *CourseDB) GetCourseInfo(course *entity.Course) error {
	return r.db.
		Preload(clause.Associations).
		First(course).
		Error
}

func (r *CourseDB) GetCourseVideo(courseID uuid.UUID, courseVideo *[]entity.CourseVideo) error {
	return r.db.
		Model(&entity.CourseVideo{}).
		Where("course_id = ?", courseID).
		Find(courseVideo).
		Error
}

func (r *CourseDB) GetCourseModule(courseID uuid.UUID, courseModule *[]entity.CourseModule) error {
	return r.db.
		Preload(clause.Associations).
		Where("course_id = ?", courseID).
		Find(courseModule).
		Error
}

func (r *CourseDB) GetCourseEnrollmentCount(courseID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.
		Model(entity.UserCourse{}).
		Where("course_id = ?", courseID).
		Count(&count).
		Error

	return count, err
}

func (r *CourseDB) UpdateCourseEnrollment(userCourse *entity.UserCourse) error {
	return r.db.
		Model(&entity.UserCourse{}).
		Create(userCourse).
		Error
}
