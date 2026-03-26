// Package ai connects and tests AI from multiple providers
package ai

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type AIItf interface {
	NewAIInterview(ctx context.Context, newAIInterview dto.NewAIInterview) (string, string, error)
	ContinueAIInterview(ctx context.Context, continueAIInterview dto.ContinueAIInterview) (string, string, error)
	UploadCV(ctx context.Context, file *multipart.FileHeader) (string, error)
	AnalyzeCV(ctx context.Context, analyzeCV dto.AnalyzeCV) (string, error)
}

type AI struct {
	OpenAI *openai.Client
	env    *env.Env
}

func New(env *env.Env) *AI {
	openAI := openai.NewClient(
		option.WithAPIKey(env.OpenAIAPIKey),
	)

	AI := AI{
		OpenAI: &openAI,
		env:    env,
	}

	return &AI
}

func Test(a *AI) {
	log.Println("testing openai connection")

	_, err := a.OpenAI.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say \"a\""),
		},
		Model: a.env.OpenAIAllowedModel,
	})
	if err != nil {
		log.Panic("openai connection failed")
	}

	log.Println("openai connection success")
}

func (a *AI) NewAIInterview(ctx context.Context, newAIInterview dto.NewAIInterview) (string, string, error) {
	params := responses.ResponseNewParams{
		Model:        a.env.OpenAIAllowedModel,
		Instructions: openai.String("You are a professional interviewer. Ask one question at a time, build on the user's previous answers, maintain a neutral tone, progress from broad to specific questions, and provide constructive feedback only when the interview is complete."),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(newAIInterview.Input),
		},
	}

	response, err := a.OpenAI.Responses.New(ctx, params)

	return response.ID, response.OutputText(), err
}

func (a *AI) ContinueAIInterview(ctx context.Context, continueAIInterview dto.ContinueAIInterview) (string, string, error) {
	params := responses.ResponseNewParams{
		Model:              a.env.OpenAIAllowedModel,
		PreviousResponseID: openai.String(continueAIInterview.PreviousResponseID),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(continueAIInterview.Input),
		},
	}

	response, err := a.OpenAI.Responses.New(ctx, params)

	return response.ID, response.OutputText(), err
}

func (a *AI) UploadCV(ctx context.Context, file *multipart.FileHeader) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", err
	}

	inputfile := openai.File(fileReader, file.Filename, "application/pdf")

	openAIFileExpirySeconds := a.env.OpenAIFileExpirySeconds
	if openAIFileExpirySeconds < 3600 {
		openAIFileExpirySeconds = 3600
	}

	storedFile, err := a.OpenAI.Files.New(
		ctx,
		openai.FileNewParams{
			File:    inputfile,
			Purpose: openai.FilePurposeUserData,
			ExpiresAfter: openai.FileNewParamsExpiresAfter{
				Anchor:  "created_at",
				Seconds: openAIFileExpirySeconds,
			},
		},
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return storedFile.ID, nil
}

func (a *AI) AnalyzeCV(ctx context.Context, analyzeCV dto.AnalyzeCV) (string, error) {
	params := responses.ResponseNewParams{
		Model:        a.env.OpenAIAllowedModel,
		Instructions: openai.String("Analyze the provided CV strictly. Your response must only contain the analysis result. Do not offer to create, rewrite, modify, or improve the CV. Do not suggest next steps. Do not ask clarifying questions. Do not add any closing remarks or offers for further assistance."),
	}

	params.Input = responses.ResponseNewParamsInputUnion{
		OfInputItemList: responses.ResponseInputParam{
			responses.ResponseInputItemParamOfMessage(
				responses.ResponseInputMessageContentListParam{
					responses.ResponseInputContentUnionParam{
						OfInputFile: &responses.ResponseInputFileParam{
							FileID: openai.String(analyzeCV.FileID),
							Type:   "input_file",
						},
					},
					responses.ResponseInputContentUnionParam{
						OfInputText: &responses.ResponseInputTextParam{
							Text: "Analyze my CV.",
							Type: "input_text",
						},
					},
					responses.ResponseInputContentUnionParam{
						OfInputText: &responses.ResponseInputTextParam{
							Text: fmt.Sprintf("%+v", analyzeCV),
							Type: "input_text",
						},
					},
				},
				"user",
			),
		},
	}

	res, err := a.OpenAI.Responses.New(ctx, params)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.OutputText(), nil
}
