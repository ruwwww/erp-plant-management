package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MediaServiceImpl struct {
	mediaRepo     repository.MediaRepository
	mediaLinkRepo repository.MediaLinkRepository
	minioClient   *minio.Client
	bucketName    string
}

func NewMediaService(mediaRepo repository.MediaRepository, mediaLinkRepo repository.MediaLinkRepository) (MediaService, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "media-assets"
	}
	secure := os.Getenv("MINIO_SECURE") == "true"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MediaServiceImpl{
		mediaRepo:     mediaRepo,
		mediaLinkRepo: mediaLinkRepo,
		minioClient:   minioClient,
		bucketName:    bucketName,
	}, nil
}

func (s *MediaServiceImpl) InitiateUpload(ctx context.Context, filename string, mimeType string, sizeBytes int64) (*domain.MediaAsset, string, error) {
	// Generate UUID for unique path
	fileUUID := uuid.New().String()
	ext := filepath.Ext(filename)
	path := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), fileUUID, ext)

	// Create MediaAsset record
	asset := &domain.MediaAsset{
		UUID:      fileUUID,
		Disk:      "s3_main",
		Path:      path,
		Filename:  filename,
		MimeType:  mimeType,
		SizeBytes: sizeBytes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.mediaRepo.Create(ctx, asset); err != nil {
		return nil, "", fmt.Errorf("failed to create media asset: %w", err)
	}

	// Generate signed PUT URL
	signedURL, err := s.minioClient.PresignedPutObject(ctx, s.bucketName, path, time.Hour*24) // 24 hours expiry
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return asset, signedURL.String(), nil
}

func (s *MediaServiceImpl) LinkMedia(ctx context.Context, mediaIDs []int, entityType string, entityID int, zone string) error {
	for i, mediaID := range mediaIDs {
		link := &domain.MediaLink{
			MediaID:    mediaID,
			EntityType: entityType,
			EntityID:   entityID,
			Zone:       zone,
			SortOrder:  i,
			CreatedAt:  time.Now(),
		}
		if err := s.mediaLinkRepo.Create(ctx, link); err != nil {
			return fmt.Errorf("failed to link media %d: %w", mediaID, err)
		}
	}
	return nil
}

func (s *MediaServiceImpl) UnlinkMedia(ctx context.Context, mediaID int, entityType string, entityID int) error {
	// Find and delete the specific link
	links, err := s.mediaLinkRepo.Find(ctx, "media_id = ? AND entity_type = ? AND entity_id = ?", mediaID, entityType, entityID)
	if err != nil {
		return fmt.Errorf("failed to find media link: %w", err)
	}
	if len(links) == 0 {
		return fmt.Errorf("media link not found")
	}

	return s.mediaLinkRepo.Delete(ctx, links[0].ID)
}

func (s *MediaServiceImpl) GetSignedURL(ctx context.Context, mediaID int) (string, error) {
	asset, err := s.mediaRepo.FindByID(ctx, mediaID)
	if err != nil {
		return "", fmt.Errorf("failed to find media asset: %w", err)
	}

	// Generate signed GET URL
	signedURL, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, asset.Path, time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed GET URL: %w", err)
	}

	return signedURL.String(), nil
}
