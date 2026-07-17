package infrastructure

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/common/apperror"
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(ctx context.Context, cfg config.AWSConfig) (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	customCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, apperror.ErrInternal.WithRootErr(err).WithMessage("Failed to load AWS configuration")
	}

	client := s3.NewFromConfig(customCfg)
	return client, nil
}
