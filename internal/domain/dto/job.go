// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResponseSearchJob struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	Location               string    `json:"location"`
	MinimumExperienceYears uint      `json:"minimum_experience_years"`
	Seniority              string    `json:"seniority"`
	Contract               string    `json:"contact"`
	SalaryMonthRangeLow    uint32    `json:"salary_month_range_low"`
	SalaryMonthRangeHigh   uint32    `json:"salary_month_range_high"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type ResponseGetJobInfo struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	Location               string    `json:"location"`
	MinimumExperienceYears uint      `json:"minimum_experience_years"`
	Seniority              string    `json:"seniority"`
	Contract               string    `json:"contact"`
	SalaryMonthRangeLow    uint32    `json:"salary_month_range_low"`
	SalaryMonthRangeHigh   uint32    `json:"salary_month_range_high"`
	JobDetail              struct {
		TeamSize      uint      `json:"team_size"`
		HiringManager string    `json:"hiring_manager"`
		WorkSetup     string    `json:"work_setup"`
		Posted        time.Time `json:"posted_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Deadline      time.Time `json:"deadline"`
	} `json:"job_detail"`
	JobTag []struct {
		Tag string `json:"tag"`
	} `json:"tag"`
	JobTask []struct {
		Task string `json:"task"`
	} `json:"job_task"`
	JobRequirement []struct {
		Requirement string `json:"task"`
	} `json:"job_requirement"`
	JobBenefit []struct {
		Benefit string `json:"benefit"`
	} `json:"job_benefit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetJobList struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	Location               string    `json:"location"`
	MinimumExperienceYears uint      `json:"minimum_experience_years"`
	Seniority              string    `json:"seniority"`
	Contract               string    `json:"contact"`
	SalaryMonthRangeLow    uint32    `json:"salary_month_range_low"`
	SalaryMonthRangeHigh   uint32    `json:"salary_month_range_high"`
	JobTag                 []struct {
		Tag string `json:"tag"`
	} `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetJobTag struct {
	Tag string `json:"tag"`
}

type ResponseGetCompanyInfo struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	ProfilePicture string    `json:"profile_picture"`
}
