// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"github.com/google/uuid"
)

type AnalyzeCV struct {
	ID                       uuid.UUID `json:"id"`
	UserID                   uuid.UUID `json:"user_id"`
	FileID                   string    `json:"file_id" validate:"required"`
	JobTitle                 string    `json:"job_title" validate:"required"`
	TargetCompany            string    `json:"target_company" validate:"required"`
	Industry                 string    `json:"industry" validate:"required"`
	WorkExperience           string    `json:"work_experience" validate:"required"`
	HighestEducation         string    `json:"highest_education" validate:"required"`
	EmploymentStatus         string    `json:"employment_status" validate:"required"`
	PrimarySkill             string    `json:"primary_skill" validate:"required"`
	Tools                    string    `json:"tools" validate:"required"`
	SpokenAndWrittenLanguage string    `json:"spoken_and_written_language" validate:"required"`
	PrimaryAnalysisGoals     string    `json:"primary_analysis_goals" validate:"required"`
	JobApplicationsSent      string    `json:"job_applications_sent" validate:"required"`
	BiggestConcern           string    `json:"biggest_concern" validate:"required"`
	AdditionalRequest        string    `json:"addititional_request"`
	ComprehensiveModel       bool      `json:"comprehensive_model" validate:"omitempty,boolean"`
}

type NewAIInterview struct {
	JobRole       string `json:"job_role" validate:"required"`
	InterviewType string `json:"interview_type" validate:"required"`
	Difficulty    string `json:"difficulty" validate:"required"`
}

type ContinueAIInterview struct {
	PreviousResponseID string `json:"previous_response_id" validate:"required"`
	Input              string `json:"input" validate:"required,min=1,max=8192"`
}

type ResponseAnalyzeCV struct {
	ID       uuid.UUID `json:"id"`
	FileID   string    `json:"file_id"`
	Response string    `json:"response"`
}

type ResponseUploadCV struct {
	FileID string `json:"file_id"`
}

type ResponseAIInterview struct {
	PreviousResponseID string `json:"previous_response_id"`
	Response           string `json:"response"`
}

type ResponseTranscribe struct {
	Response string `json:"response"`
}
