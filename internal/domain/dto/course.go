// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResponseGetCourseList struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CoverImage string    `json:"cover_image"`
	Price      uint32    `json:"price"`
}

type ResponseGetCourseInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	Price       uint32    `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResponseSearchCourse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CoverImage string    `json:"cover_image"`
	Price      uint32    `json:"price"`
}
