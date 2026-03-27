// Package ai connects and tests AI from multiple providers
package ai

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type AIItf interface {
	NewAIInterview(ctx context.Context, newAIInterview dto.NewAIInterview) (string, string, error)
	ContinueAIInterview(ctx context.Context, continueAIInterview dto.ContinueAIInterview) (string, string, error)
	Transcribe(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadCV(ctx context.Context, file *multipart.FileHeader) (string, error)
	AnalyzeCV(ctx context.Context, analyzeCV dto.AnalyzeCV) (string, error)
}

type AI struct {
	OpenAI         *openai.Client
	AIInstructions AIInstructions
	env            *env.Env
}

type AIInstructions struct {
	Interview string
	AnalyzeCV string
}

func New(env *env.Env) *AI {
	openAI := openai.NewClient(
		option.WithAPIKey(env.OpenAIAPIKey),
	)

	AIInstructions := AIInstructions{
		Interview: `
					You are conducting a professional interview. You are the interviewer. The user is the candidate.

					IMPORTANT: You may only send one question per message. If you find yourself writing more than one question, delete all but the most important one before sending.

					Your only job is to ask questions and listen. Do not offer opinions, tips, explanations, or encouragement during the interview. Do not break character under any circumstances.

					Rules:
					- Start the interview immediately by introducing yourself briefly and asking your first question
					- Ask exactly one question per message, no exceptions
					- Never use bullet points, numbered lists, or multiple questions in a single message, even as follow-ups or examples
					- Base each follow-up on what the candidate just said
					- Move from broad, open-ended questions toward specific, probing ones as the interview progresses
					- If the candidate goes off-topic, redirect them with a short, neutral phrase and continue
					- Do not say things like "great answer", "interesting", or "that's a good point"
					- Do not reveal that you are an AI or a language model

					End condition:
					- When the interview is complete when you judge the topic is covered, say "That concludes our interview" and then provide structured feedback covering: strengths, areas for improvement, and an overall impression.
			`,
		AnalyzeCV: "Analyze the provided CV strictly. Your response must only contain the analysis result. Do not offer to create, rewrite, modify, or improve the CV. Do not suggest next steps. Do not ask clarifying questions. Do not add any closing remarks or offers for further assistance.",
	}

	AI := AI{
		OpenAI:         &openAI,
		AIInstructions: AIInstructions,
		env:            env,
	}

	return &AI
}

func Test(a *AI) {
	log.Println("testing openai connection")

	_, err := a.OpenAI.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say \"a\""),
		},
		Model: a.env.OpenAITextModel,
	})
	if err != nil {
		log.Panic("openai connection failed")
	}

	log.Println("openai connection success")
}

func (a *AI) NewAIInterview(ctx context.Context, newAIInterview dto.NewAIInterview) (string, string, error) {
	params := responses.ResponseNewParams{
		Model:        a.env.OpenAITextModel,
		Instructions: openai.String(a.AIInstructions.Interview),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(fmt.Sprintf("%+v", newAIInterview)),
		},
	}

	response, err := a.OpenAI.Responses.New(ctx, params)

	return response.ID, response.OutputText(), err
}

func (a *AI) ContinueAIInterview(ctx context.Context, continueAIInterview dto.ContinueAIInterview) (string, string, error) {
	params := responses.ResponseNewParams{
		Model:              a.env.OpenAITextModel,
		PreviousResponseID: openai.String(continueAIInterview.PreviousResponseID),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(continueAIInterview.Input),
		},
	}

	response, err := a.OpenAI.Responses.New(ctx, params)

	return response.ID, response.OutputText(), err
}

func (a *AI) Transcribe(ctx context.Context, file *multipart.FileHeader) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", err
	}

	mimeType, err := mimetype.DetectReader(fileReader)
	if err != nil || !mimeType.Is("audio/mpeg") {
		return "", fiber.ErrBadRequest
	}

	_, err = fileReader.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	inputFile := openai.File(fileReader, file.Filename, mimeType.String())

	params := openai.AudioTranscriptionNewParams{
		Model: a.env.OpenAITranscribeModel,
		File:  inputFile,
	}

	transcription, err := a.OpenAI.Audio.Transcriptions.New(ctx, params)
	if err != nil {
		return "", err
	}

	return transcription.Text, nil
}

func (a *AI) UploadCV(ctx context.Context, file *multipart.FileHeader) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", err
	}

	mimeType, err := mimetype.DetectReader(fileReader)
	if err != nil || !mimeType.Is("application/pdf") {
		return "", fiber.ErrBadRequest
	}

	_, err = fileReader.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	inputFile := openai.File(fileReader, file.Filename, mimeType.String())

	openAIFileExpirySeconds := a.env.OpenAIFileExpirySeconds
	if openAIFileExpirySeconds < 3600 {
		openAIFileExpirySeconds = 3600
	}

	storedFile, err := a.OpenAI.Files.New(
		ctx,
		openai.FileNewParams{
			File:    inputFile,
			Purpose: openai.FilePurposeUserData,
			ExpiresAfter: openai.FileNewParamsExpiresAfter{
				Anchor:  "created_at",
				Seconds: openAIFileExpirySeconds,
			},
		},
	)
	if err != nil {
		return "", err
	}

	return storedFile.ID, nil
}

func (a *AI) AnalyzeCV(ctx context.Context, analyzeCV dto.AnalyzeCV) (string, error) {
	params := responses.ResponseNewParams{
		Model:        a.env.OpenAITextModel,
		Instructions: openai.String(a.AIInstructions.AnalyzeCV),
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
		return "", err
	}

	return res.OutputText(), nil
}
