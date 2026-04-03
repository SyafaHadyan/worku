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
	GetCourseCategory() ([]dto.ResponseGetCourseCategory, error)
	GetCourseList(offset int, limit int) ([]dto.ResponseGetCourseList, error)
	GetCourseListByCategory(categoryID uuid.UUID, offset int, limit int) ([]dto.ResponseGetCourseList, error)
	GetCourseInfo(userID uuid.UUID, courseID uuid.UUID) (dto.ResponseGetCourseInfo, error)
	SearchCourse(offset int, limit int, query string) ([]dto.ResponseSearchCourse, error)
	GetCourseVideo(courseID uuid.UUID) ([]dto.ResponseGetCourseVideo, error)
	GetCourseModule(courseID uuid.UUID) ([]dto.ResponseGetCourseModule, error)
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

func (u *CourseUseCase) GetCourseCategory() ([]dto.ResponseGetCourseCategory, error) {
	var courseCategory []entity.CourseCategory

	err := u.courseRepo.GetCourseCategory(&courseCategory)
	if len(courseCategory) == 0 {
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	courseCategoryList := make([]dto.ResponseGetCourseCategory, len(courseCategory))

	for i, courseCategoryItem := range courseCategory {
		courseCategoryList[i] = courseCategoryItem.ParseToDTOResponseGetCourseCategory()
	}

	return courseCategoryList, nil
}

func (u *CourseUseCase) GetCourseList(offset int, limit int) ([]dto.ResponseGetCourseList, error) {
	var course []entity.Course

	offset = offset * limit

	err := u.courseRepo.GetCourseList(&offset, &limit, &course)
	if len(course) == 0 {
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	courseList := make([]dto.ResponseGetCourseList, len(course))

	for i, courseItem := range course {
		courseList[i] = courseItem.ParseToDTOResponseGetCourseList()
	}

	return courseList, nil
}

func (u *CourseUseCase) GetCourseListByCategory(categoryID uuid.UUID, offset int, limit int) ([]dto.ResponseGetCourseList, error) {
	var course []entity.Course

	offset = offset * limit

	err := u.courseRepo.GetCourseListByCategory(categoryID, &offset, &limit, &course)
	if len(course) == 0 {
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	courseList := make([]dto.ResponseGetCourseList, len(course))

	for i, courseItem := range course {
		courseList[i] = courseItem.ParseToDTOResponseGetCourseList()
	}

	return courseList, nil
}

func (u *CourseUseCase) GetCourseInfo(userID uuid.UUID, courseID uuid.UUID) (dto.ResponseGetCourseInfo, error) {
	var courseParsed dto.ResponseGetCourseInfo

	course := entity.Course{
		ID: courseID,
	}

	userCourse := entity.UserCourse{
		UserID:   userID,
		CourseID: courseID,
	}

	redisKey := fmt.Sprintf("course:%s", courseID.String())

	result, err := u.redis.Get(redisKey)
	if err == nil && result != "" {
		var out dto.ResponseGetCourseInfo

		err := json.Unmarshal([]byte(result), &out)
		if err != nil {
			log.Println(err)
		}

		return out, nil
	}

	err = u.courseRepo.GetCourseInfo(&course)
	if err != nil {
		return dto.ResponseGetCourseInfo{}, err
	}

	go func() {
		err := u.courseRepo.UpdateCourseEnrollment(&userCourse)
		if err != nil {
			log.Println(err)
		}
	}()

	enrollmentCount, err := u.courseRepo.GetCourseEnrollmentCount(courseID)
	if err != nil {
		log.Println(err)
	}

	courseParsed = course.ParseToDTOResponseGetCourseInfo()
	courseParsed.EnrollmentCount = enrollmentCount

	go func() {
		newData, err := json.Marshal(courseParsed)
		if err != nil {
			log.Println(err)
		}

		u.redis.Set(redisKey, string(newData))
	}()

	return courseParsed, nil
}

func (u *CourseUseCase) SearchCourse(offset int, limit int, query string) ([]dto.ResponseSearchCourse, error) {
	var course []entity.Course

	offset = offset * limit

	err := u.courseRepo.SearchCourse(&offset, &limit, &query, &course)
	if err != nil {
		return nil, err
	}

	courseList := make([]dto.ResponseSearchCourse, len(course))

	for i, courseItem := range course {
		courseList[i] = courseItem.ParseToDTOResponseSearchCourse()
	}

	return courseList, nil
}

func (u *CourseUseCase) GetCourseVideo(courseID uuid.UUID) ([]dto.ResponseGetCourseVideo, error) {
	var courseVideo []entity.CourseVideo

	err := u.courseRepo.GetCourseVideo(courseID, &courseVideo)
	if err != nil {
		return nil, err
	}

	courseVideoList := make([]dto.ResponseGetCourseVideo, len(courseVideo))

	for i, courseVideoItem := range courseVideo {
		courseVideoList[i] = courseVideoItem.ParseToDTOResponseGetCourseVideo()
	}

	return courseVideoList, nil
}

func (u *CourseUseCase) GetCourseModule(courseID uuid.UUID) ([]dto.ResponseGetCourseModule, error) {
	var courseModule []entity.CourseModule

	err := u.courseRepo.GetCourseModule(courseID, &courseModule)
	if err != nil {
		return nil, err
	}

	courseModuleList := make([]dto.ResponseGetCourseModule, len(courseModule))

	for i, courseModuleItem := range courseModule {
		courseModuleList[i] = courseModuleItem.ParseToDTOResponseGetCourseModule()
	}

	return courseModuleList, nil
}
