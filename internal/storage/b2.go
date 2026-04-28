package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "github.com/AdelmoMJunior/GoACBr/internal/config"
)

// B2Storage implementation of Provider for Backblaze B2 (S3 compatible).
type B2Storage struct {
	client       *s3.Client
	bucketName   string
	presignCli   *s3.PresignClient
	publicCDNURL string
}

// NewB2Storage initializes a new B2 storage provider.
func NewB2Storage(cfg appconfig.B2Config) (*B2Storage, error) {
	if cfg.KeyID == "" || cfg.AppKey == "" || cfg.Endpoint == "" {
		slog.Warn("B2 storage configuration is incomplete. Storage features will fail.")
		return nil, fmt.Errorf("b2 configuration incomplete")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.KeyID, cfg.AppKey, "")),
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load B2 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	presignCli := s3.NewPresignClient(client)

	slog.Info("Initialized B2 Storage Provider", "bucket", cfg.BucketName, "endpoint", cfg.Endpoint)

	return &B2Storage{
		client:       client,
		bucketName:   cfg.BucketName,
		presignCli:   presignCli,
		publicCDNURL: cfg.PublicCDNURL,
	}, nil
}

func (s *B2Storage) Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		slog.Error("B2 upload error", "key", key, "error", err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return key, nil
}

func (s *B2Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("B2 download error", "key", key, "error", err)
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return output.Body, nil
}

func (s *B2Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("B2 delete error", "key", key, "error", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *B2Storage) GetURL(ctx context.Context, key string) (string, error) {
	if s.publicCDNURL != "" {
		return fmt.Sprintf("%s/%s/%s", s.publicCDNURL, s.bucketName, key), nil
	}

	req, err := s.presignCli.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 1 * time.Hour
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}
