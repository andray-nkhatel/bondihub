package services

import (
	"bondihub/config"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService handles image uploads to Cloudinary
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new Cloudinary service instance
// It supports both CLOUDINARY_URL format (cloudinary://key:secret@cloud) and individual parameters
func NewCloudinaryService() (*CloudinaryService, error) {
	var cld *cloudinary.Cloudinary
	var err error

	// Try to use CLOUDINARY_URL if available (preferred method)
	cloudinaryURL := config.AppConfig.CloudinaryURL
	if cloudinaryURL != "" {
		cld, err = cloudinary.NewFromURL(cloudinaryURL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cloudinary from URL: %w", err)
		}
	} else {
		// Fall back to individual parameters
		cld, err = cloudinary.NewFromParams(
			config.AppConfig.CloudinaryCloud,
			config.AppConfig.CloudinaryKey,
			config.AppConfig.CloudinarySecret,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cloudinary from params: %w", err)
		}
	}

	return &CloudinaryService{cld: cld}, nil
}

// UploadImage uploads an image file to Cloudinary
func (cs *CloudinaryService) UploadImage(ctx context.Context, file multipart.File, folder string) (*uploader.UploadResult, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert bytes to Reader for Cloudinary
	fileReader := bytes.NewReader(fileBytes)

	// Upload to Cloudinary
	result, err := cs.cld.Upload.Upload(
		ctx,
		fileReader,
		uploader.UploadParams{
			Folder:         folder,
			ResourceType:   "image",
			Transformation: "f_auto,q_auto",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return result, nil
}

// UploadImageFromBytes uploads an image from bytes to Cloudinary
func (cs *CloudinaryService) UploadImageFromBytes(ctx context.Context, fileBytes []byte, folder string) (*uploader.UploadResult, error) {
	// Convert bytes to Reader for Cloudinary
	fileReader := bytes.NewReader(fileBytes)

	result, err := cs.cld.Upload.Upload(
		ctx,
		fileReader,
		uploader.UploadParams{
			Folder:         folder,
			ResourceType:   "image",
			Transformation: "f_auto,q_auto",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return result, nil
}

// DeleteImage deletes an image from Cloudinary
func (cs *CloudinaryService) DeleteImage(ctx context.Context, publicID string) (*uploader.DestroyResult, error) {
	result, err := cs.cld.Upload.Destroy(
		ctx,
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	return result, nil
}

// GetImageURL generates a Cloudinary URL for an image
func (cs *CloudinaryService) GetImageURL(publicID string, transformations string) string {
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s/%s",
		config.AppConfig.CloudinaryCloud, transformations, publicID)
}

// GetOptimizedImageURL generates an optimized image URL
func (cs *CloudinaryService) GetOptimizedImageURL(publicID string, width, height int) string {
	transformations := fmt.Sprintf("f_auto,q_auto,w_%d,h_%d,c_fill", width, height)
	return cs.GetImageURL(publicID, transformations)
}
