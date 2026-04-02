// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResponseGetCourseCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
}

type ResponseGetCourseList struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	CoverImage string    `json:"cover_image"`
}

type ResponseSearchCourse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	CoverImage string    `json:"cover_image"`
}

type ResponseGetCourseInfo struct {
	ID              uuid.UUID `json:"id"`
	CategoryID      uuid.UUID `json:"category_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	LessonCount     int       `json:"lesson_count"`
	ModuleCount     int       `json:"module_count"`
	EnrollmentCount int64     `json:"enrollment_count"`
	CoverImage      string    `json:"cover_image"`
	CourseVideo     []struct {
		ID       uuid.UUID `json:"id"`
		VideoURL string    `json:"video_url"`
	} `json:"course_video"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetCourseVideo struct {
	ID       uuid.UUID `json:"video_id"`
	CourseID uuid.UUID `json:"course_id"`
	VideoURL string    `json:"video_url"`
}
