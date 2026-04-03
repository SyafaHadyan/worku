// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/app/job/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	redisitf "github.com/SyafaHadyan/worku/internal/infra/redis"
	s3itf "github.com/SyafaHadyan/worku/internal/infra/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobUseCaseItf interface {
	GetJobInfo(jobID uuid.UUID) (dto.ResponseGetJobInfo, error)
	GetJobList(offset int, limit int) ([]dto.ResponseGetJobList, error)
	GetCompanyInfo(companyID uuid.UUID) (dto.ResponseGetCompanyInfo, error)
}

type JobUseCase struct {
	jobRepo      repository.JobDBItf
	redis        redisitf.RedisItf
	redisContext context.Context
	s3           s3itf.S3Itf
	s3Context    context.Context
	env          *env.Env
}

func NewJobUseCase(
	jobRepo repository.JobDBItf, redis redisitf.RedisItf,
	s3 s3itf.S3Itf, env *env.Env,
) JobUseCaseItf {
	return &JobUseCase{
		jobRepo:      jobRepo,
		redis:        redis,
		redisContext: context.Background(),
		s3:           s3,
		s3Context:    context.Background(),
		env:          env,
	}
}

func (u *JobUseCase) GetJobInfo(jobID uuid.UUID) (dto.ResponseGetJobInfo, error) {
	job := entity.Job{
		ID: jobID,
	}

	redisKey := fmt.Sprintf("job:%s", jobID.String())

	result, err := u.redis.Get(redisKey)
	if err == nil && result != "" {
		var out entity.Job

		err := json.Unmarshal([]byte(result), &out)
		if err != nil {
			log.Println(err)
		}

		return out.ParseToDTOResponseGetJobInfo(), nil
	}

	err = u.jobRepo.GetJobInfo(&job)
	if err != nil {
		return dto.ResponseGetJobInfo{}, err
	}

	go func() {
		newData, err := json.Marshal(job)
		if err != nil {
			log.Println(err)
		}
		u.redis.Set(redisKey, string(newData))
	}()

	return job.ParseToDTOResponseGetJobInfo(), nil
}

func (u *JobUseCase) GetJobList(offset int, limit int) ([]dto.ResponseGetJobList, error) {
	var job []entity.Job

	offset = offset * limit

	err := u.jobRepo.GetJobList(&offset, &limit, &job)
	if len(job) == 0 {
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	jobList := make([]dto.ResponseGetJobList, len(job))

	for i, jobItem := range job {
		jobList[i] = jobItem.ParseToDTOResponseGetJobList()
	}

	return jobList, nil
}

func (u *JobUseCase) GetCompanyInfo(companyID uuid.UUID) (dto.ResponseGetCompanyInfo, error) {
	company := entity.Company{
		ID: companyID,
	}

	redisKey := fmt.Sprintf("company:%s", companyID.String())

	result, err := u.redis.Get(redisKey)
	if err == nil && result != "" {
		var out entity.Company

		err := json.Unmarshal([]byte(result), &out)
		if err != nil {
			log.Println(err)
		}

		return out.ParseToDTOResponseGetCompanyInfo(), nil
	}

	err = u.jobRepo.GetCompanyInfo(&company)
	if err != nil {
		return dto.ResponseGetCompanyInfo{}, err
	}

	go func() {
		newData, err := json.Marshal(company)
		if err != nil {
			log.Println(err)
		}
		u.redis.Set(redisKey, string(newData))
	}()

	return company.ParseToDTOResponseGetCompanyInfo(), nil
}
