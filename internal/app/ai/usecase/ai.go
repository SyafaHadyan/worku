// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/SyafaHadyan/worku/internal/app/ai/repository"
	"github.com/SyafaHadyan/worku/internal/constants"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	aiitf "github.com/SyafaHadyan/worku/internal/infra/ai"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	s3itf "github.com/SyafaHadyan/worku/internal/infra/s3"
	"github.com/google/uuid"
)

type AIUseCaseItf interface {
	NewAIInterview(newAIInterview dto.NewAIInterview) (dto.ResponseAIInterview, error)
	ContinueAIInterview(continueAIInterview dto.ContinueAIInterview) (dto.ResponseAIInterview, error)
	Transcribe(userID uuid.UUID, file *multipart.FileHeader) (dto.ResponseTranscribe, error)
	UploadCV(userID uuid.UUID, file multipart.FileHeader) (dto.ResponseUploadCV, error)
	AnalyzeCV(analyzeCV dto.AnalyzeCV) (dto.ResponseAnalyzeCV, error)
}

type AIUseCase struct {
	aiRepo    repository.AIDBItf
	ai        aiitf.AIItf
	aiContext context.Context
	s3        s3itf.S3Itf
	s3Context context.Context
	env       *env.Env
}

func NewAIUseCase(
	aiRepo repository.AIDBItf, ai aiitf.AIItf,
	s3 s3itf.S3Itf, env *env.Env,
) AIUseCaseItf {
	return &AIUseCase{
		aiRepo:    aiRepo,
		ai:        ai,
		aiContext: context.Background(),
		s3:        s3,
		s3Context: context.Background(),
		env:       env,
	}
}

func (u *AIUseCase) NewAIInterview(newAIInterview dto.NewAIInterview) (dto.ResponseAIInterview, error) {
	previousResponseID, response, err := u.ai.NewAIInterview(u.aiContext, newAIInterview)

	return dto.ResponseAIInterview{
		PreviousResponseID: previousResponseID,
		Response:           response,
	}, err
}

func (u *AIUseCase) ContinueAIInterview(continueAIInterview dto.ContinueAIInterview) (dto.ResponseAIInterview, error) {
	previousResponseID, response, err := u.ai.ContinueAIInterview(u.aiContext, continueAIInterview)

	return dto.ResponseAIInterview{
		PreviousResponseID: previousResponseID,
		Response:           response,
	}, err
}

func (u *AIUseCase) Transcribe(userID uuid.UUID, file *multipart.FileHeader) (dto.ResponseTranscribe, error) {
	response, err := u.ai.Transcribe(u.aiContext, file)

	go func() {
		fileContent, err := file.Open()
		if err != nil {
			return
		}

		byteContainer, err := io.ReadAll(fileContent)
		if err != nil {
			return
		}

		objectKey := fmt.Sprintf(
			"%s/%s/%s",
			constants.UserVoiceDirectory,
			userID.String(),
			uuid.New().String(),
		)

		err = u.s3.Upload(u.s3Context, objectKey, byteContainer)
		if err != nil {
			log.Println(err)
		}
	}()

	return dto.ResponseTranscribe{
		Response: response,
	}, err
}

func (u *AIUseCase) UploadCV(userID uuid.UUID, file multipart.FileHeader) (dto.ResponseUploadCV, error) {
	fileID, err := u.ai.UploadCV(u.aiContext, &file)
	if err != nil {
		return dto.ResponseUploadCV{}, err
	}

	responseUploadCV := dto.ResponseUploadCV{
		FileID: fileID,
	}

	go func() {
		fileContent, err := file.Open()
		if err != nil {
			return
		}

		byteContainer, err := io.ReadAll(fileContent)
		if err != nil {
			return
		}

		objectKey := fmt.Sprintf(
			"%s/%s/%s",
			constants.CVDirectory,
			userID.String(),
			fileID,
		)

		err = u.s3.Upload(u.s3Context, objectKey, byteContainer)
		if err != nil {
			log.Println(err)
		}
	}()

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
