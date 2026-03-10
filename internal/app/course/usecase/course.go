// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/app/course/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	redisitf "github.com/SyafaHadyan/worku/internal/infra/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseUseCaseItf interface {
	GetCourseList(offset int, limit int) ([]dto.ResponseGetCourseList, error)
	GetCourseInfo(courseID uuid.UUID) (dto.ResponseGetCourseInfo, error)
}

type CourseUseCase struct {
	courseRepo   repository.CourseDBItf
	redis        redisitf.RedisItf
	redisContext context.Context
}

func NewCourseUseCase(
	courseRepo repository.CourseDBItf, redis redisitf.RedisItf,
) CourseUseCaseItf {
	return &CourseUseCase{
		courseRepo:   courseRepo,
		redis:        redis,
		redisContext: context.Background(),
	}
}

func (u *CourseUseCase) GetCourseList(offset int, limit int) ([]dto.ResponseGetCourseList, error) {
	var course []entity.Course

	offset = offset * limit

	err := u.courseRepo.GetCourseList(offset, limit, &course)
	if err != nil {
		return nil, err
	}

	courseList := make([]dto.ResponseGetCourseList, len(course))

	for i, courseItem := range course {
		courseList[i] = courseItem.ParseToDTOResponseGetCourseList()
	}

	return courseList, nil
}

func (u *CourseUseCase) GetCourseInfo(courseID uuid.UUID) (dto.ResponseGetCourseInfo, error) {
	var count int64

	course := entity.Course{
		ID: courseID,
	}

	redisKey := fmt.Sprintf("course:%s", courseID.String())

	result, err := u.redis.Get(redisKey)
	if err == nil && result != "" {
		var out entity.Course

		err := json.Unmarshal([]byte(result), &out)
		if err != nil {
			log.Println(err)
		}

		return out.ParseToDTOResponseGetCourseInfo(), nil
	}

	err = u.courseRepo.GetCourseInfo(&count, &course)
	if err != nil {
		return course.ParseToDTOResponseGetCourseInfo(), nil
	}

	if count == 0 {
		return dto.ResponseGetCourseInfo{}, gorm.ErrRecordNotFound
	}

	go func() {
		newData, err := json.Marshal(course)
		if err != nil {
			log.Println(err)
		}

		u.redis.Set(redisKey, string(newData))
	}()

	return course.ParseToDTOResponseGetCourseInfo(), nil
}
