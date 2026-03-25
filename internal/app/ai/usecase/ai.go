// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"mime/multipart"

	"github.com/SyafaHadyan/worku/internal/app/ai/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	aiitf "github.com/SyafaHadyan/worku/internal/infra/ai"
	"github.com/google/uuid"
)

type AIUseCaseItf interface {
	UploadCV(file multipart.FileHeader) (dto.ResponseUploadCV, error)
	AnalyzeCV(analyzeCV dto.AnalyzeCV) (dto.ResponseAnalyzeCV, error)
}

type AIUseCase struct {
	aiRepo    repository.AIDBItf
	ai        aiitf.AIItf
	aiContext context.Context
}

func NewAIUseCase(
	aiRepo repository.AIDBItf, ai aiitf.AIItf,
) AIUseCaseItf {
	return &AIUseCase{
		aiRepo:    aiRepo,
		ai:        ai,
		aiContext: context.Background(),
	}
}

func (u *AIUseCase) UploadCV(file multipart.FileHeader) (dto.ResponseUploadCV, error) {
	fileID, err := u.ai.UploadCV(u.aiContext, &file)
	if err != nil {
		return dto.ResponseUploadCV{}, err
	}

	responseUploadCV := dto.ResponseUploadCV{
		FileID: fileID,
	}

	return responseUploadCV, err
}

func (u *AIUseCase) AnalyzeCV(analyzeCV dto.AnalyzeCV) (dto.ResponseAnalyzeCV, error) {
	analyzeCV.ID = uuid.New()

	response, err := u.ai.AnalyzeCV(u.aiContext, analyzeCV)
	if err != nil {
		return dto.ResponseAnalyzeCV{}, err
	}

	responseAnalyzeCV := entity.ResponseAnalyzeCV{
		ID:       analyzeCV.ID,
		UserID:   analyzeCV.UserID,
		FileID:   analyzeCV.FileID,
		Response: response,
	}

	err = u.aiRepo.ResponseAnalyzeCV(&responseAnalyzeCV)

	return responseAnalyzeCV.ParseToDTOResponseAnalyzeCV(), err
}
