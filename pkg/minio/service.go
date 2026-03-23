package minio

import (
	"context"
	"fmt"
	"io"
	"mkit/pkg/config"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client        *minio.Client
	logger        *logrus.Logger
	defaultBucket string
}

func New(config *config.App, logger *logrus.Logger) (*Service, error) {
	cfg := config.Minio
	minioClient, err := minio.New(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.SSLEnabled,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot init minio client: %w", err)
	}

	return &Service{
		client:        minioClient,
		logger:        logger,
		defaultBucket: cfg.BucketName,
	}, nil
}

func (s *Service) PutObjectWithBucket(
	ctx context.Context, bucketName, objectName string, data io.Reader, size int64, contentType string,
) (minio.UploadInfo, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return s.client.PutObject(ctx, bucketName, objectName, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
}

func (s *Service) PutObject(
	ctx context.Context, objectName string, data io.Reader, size int64, contentType string,
) (minio.UploadInfo, error) {
	return s.PutObjectWithBucket(ctx, s.defaultBucket, objectName, data, size, contentType)
}

func (s *Service) CheckIfObjectExistsWithBucket(
	ctx context.Context, bucketName, objectName string,
) (bool, error) {
	_, err := s.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object %s: %w", objectName, err)
	}

	return true, nil
}

func (s *Service) CheckIfObjectExists(ctx context.Context, objectName string) (bool, error) {
	return s.CheckIfObjectExistsWithBucket(ctx, s.defaultBucket, objectName)
}

func (s *Service) ListObjectNamesWithBucket(ctx context.Context, bucketName, basePath string) ([]string, error) {
	if basePath != "" && !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	objectCh := s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    basePath,
		Recursive: false,
	})

	var files []string
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		files = append(files, obj.Key)
	}

	return files, nil
}

func (s *Service) ListObjectNames(ctx context.Context, basePath string) ([]string, error) {
	return s.ListObjectNamesWithBucket(ctx, s.defaultBucket, basePath)
}

func (s *Service) GetObjectWithBucket(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (s *Service) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	return s.GetObjectWithBucket(ctx, s.defaultBucket, objectName)
}
