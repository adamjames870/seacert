package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	GetPresignedUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error)
	GetPresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type R2Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucketName string
}

type Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Endpoint        string
}

func NewR2Storage(ctx context.Context, cfg Config) (*R2Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom, // Explicitly mark as custom
		}, nil
	})

	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// IMPORTANT: Use NewFromConfig with options to ensure compatibility
	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true // This is often required for specific S3-compatible providers
	})

	presigner := s3.NewPresignClient(client)

	return &R2Storage{
		client:     client,
		presigner:  presigner,
		bucketName: cfg.BucketName,
	}, nil
}

func (s *R2Storage) GetPresignedUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error) {
	request, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("failed to presign upload URL: %w", err)
	}

	return request.URL, nil
}

func (s *R2Storage) GetPresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	request, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(s.bucketName),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("inline"), // Add this line
	}, s3.WithPresignExpires(expires))

	if err != nil {
		return "", fmt.Errorf("failed to presign download URL: %w", err)
	}

	return request.URL, nil
}

func (s *R2Storage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from R2: %w", err)
	}
	return nil
}
