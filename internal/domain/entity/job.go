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

type Job struct {
	ID                     uuid.UUID `gorm:"type:char(36);primaryKey"`
	CompanyID              uuid.UUID `gorm:"type:char(36);not null;unique"`
	Name                   string    `gorm:"type:nvarchar(128);index:idx_fulltext_search,class:FULLTEXT"`
	Location               string    `gorm:"type:nvarchar(128);index:idx_fulltext_search,class:FULLTEXT"`
	MinimumExperienceYears uint      `gorm:"type:integer unsigned"`
	Seniority              string    `gorm:"type:nvarchar(64)"`
	Contract               string    `gorm:"type:nvarchar(64)"`
	SalaryMonthRangeLow    uint32    `gorm:"type:integer unsigned"`
	SalaryMonthRangeHigh   uint32    `gorm:"type:integer unsigned"`
	JobDetail              JobDetail
	JobTag                 []JobTag
	JobTask                []JobTask
	JobRequirement         []JobRequirement
	JobBenefit             []JobBenefit
	CreatedAt              time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt              time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type JobDetail struct {
	JobID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	TeamSize      uint      `gorm:"type:integer unsigned"`
	HiringManager string    `gorm:"type:nvarchar(128)"`
	WorkSetup     string    `gorm:"type:nvarchar(64)"`
	PostedAt      time.Time `gorm:"type:timestamp;autoCreateTime"`
	Deadline      time.Time `gorm:"type:timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

type JobTag struct {
	JobID uuid.UUID `gorm:"type:char(36);primaryKey"`
	Tag   string    `gorm:"type:nvarchar(128);primaryKey"`
}

type JobTask struct {
	JobID uuid.UUID `gorm:"type:char(36);primaryKey"`
	Task  string    `gorm:"type:nvarchar(128);primaryKey"`
}

type JobRequirement struct {
	JobID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Requirement string    `gorm:"type:nvarchar(128);primaryKey"`
}

type JobBenefit struct {
	JobID   uuid.UUID `gorm:"type:char(36);primaryKey"`
	Benefit string    `gorm:"type:nvarchar(128);primaryKey"`
}

type Company struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey"`
	Name           string         `gorm:"type:nvarchar(128)"`
	ProfilePicture string         `gorm:"type:nvarchar(512)"`
	CreatedAt      time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (j *Job) ParseToDTOResponseGetJobInfo() dto.ResponseGetJobInfo {
	var response dto.ResponseGetJobInfo

	err := copier.Copy(&response, j)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (j *Job) ParseToDTOResponseGetJobList() dto.ResponseGetJobList {
	var response dto.ResponseGetJobList

	err := copier.Copy(&response, j)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (j *Job) ParseToDTOResponseSearchJob() dto.ResponseSearchJob {
	var response dto.ResponseSearchJob

	err := copier.Copy(&response, j)
	if err != nil {
		log.Println(err)
	}

	return response
}

func (j *Company) ParseToDTOResponseGetCompanyInfo() dto.ResponseGetCompanyInfo {
	var response dto.ResponseGetCompanyInfo

	err := copier.Copy(&response, j)
	if err != nil {
		log.Println(err)
	}

	return response
}
