package service

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var safeExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type StorageService struct {
	db       *gorm.DB
	s3Client *s3.Client
	bucket   string
}

func NewStorageService(db *gorm.DB, s3Client *s3.Client, bucket string) *StorageService {
	return &StorageService{db: db, s3Client: s3Client, bucket: bucket}
}

func (s *StorageService) GeneratePresignedURL(ctx context.Context, userID uuid.UUID, req dto.FileUploadRequest) (*dto.FileURLResponse, error) {
	ext, ok := safeExtensions[req.ContentType]
	if !ok {
		return nil, apperror.ErrUnsupportedFileType
	}

	if req.Size > 52422880 {
		return nil, apperror.ErrFileTooLarge
	}

	fileUUID := uuid.New()
	key := fmt.Sprintf("workspaces/assets/%s/%s%s", userID.String(), fileUUID, ext)
	presignClient := s3.NewPresignClient(s.s3Client)

	presign, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(req.ContentType),
		ContentLength: aws.Int64(req.Size), // FE should send header "Content-Length" with file size
	}, s3.WithPresignExpires(10*time.Minute))

	if err != nil {
		return nil, apperror.ErrInternal.WithRootErr(err).WithMessage("Failed to generate presigned URL")
	}

	fileRecord := &model.File{
		FileID:        fileUUID,
		FileUserID:    userID,
		FileKey:       key,
		FileStatus:    model.FileStatusPending,
		FileExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.db.WithContext(ctx).Create(&fileRecord).Error; err != nil {
		return nil, apperror.ErrInternal.WithRootErr(err).WithMessage("Failed to create file")
	}

	return &dto.FileURLResponse{
		UploadURL: presign.URL,
		FileID:    fileUUID,
		Key:       key,
		ExpiresIn: 600, // 10 minutes in seconds
	}, nil
}
