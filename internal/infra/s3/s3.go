// Package s3 connects and abstracts s3 storage implementation
package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Itf interface {
	Upload(ctx context.Context, objectKey string, object []byte) error
}

type S3 struct {
	Client *s3.Client
	env    *env.Env
}

func New(env *env.Env) *S3 {
	config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.S3AccessKeyID, env.S3AccessKeySecret, "")),
		config.WithRegion("auto"))
	if err != nil {
		log.Panic(err)
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", env.S3AccountID))
	})

	S3 := S3{
		Client: client,
		env:    env,
	}

	return &S3
}

func (s *S3) Test() {
	output, err := s.Client.ListObjectsV2(
		context.Background(),
		&s3.ListObjectsV2Input{
			Bucket: aws.String(s.env.S3BucketName),
		})
	if err != nil {
		log.Println(err)
	}

	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), *object.Size)
	}
}

func (s *S3) Upload(ctx context.Context, objectKey string, object []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.env.S3BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(object),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			errorString := errors.New(fmt.Sprintf("S3: %v\n", "file too large"))

			log.Printf(errorString.Error())

			return errorString
		}
	}

	return nil
}
