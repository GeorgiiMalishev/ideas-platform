package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageUsecaseImpl struct {
	minioClient *minio.Client
	bucketName  string
}

func NewImageUsecase(minioClient *minio.Client, bucketName string) *ImageUsecaseImpl {
	return &ImageUsecaseImpl{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (uc *ImageUsecaseImpl) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := uuid.New().String() + filepath.Ext(file.Filename)

	_, err = uc.minioClient.PutObject(ctx, uc.bucketName, fileName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", uc.bucketName, fileName), nil
}

func (uc *ImageUsecaseImpl) CreateBucket(ctx context.Context) error {
	found, err := uc.minioClient.BucketExists(ctx, uc.bucketName)
	if err != nil {
		return err
	}
	if !found {
		err = uc.minioClient.MakeBucket(ctx, uc.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *ImageUsecaseImpl) GetImage(ctx context.Context, objectName string) (*minio.Object, minio.ObjectInfo, error) {
	object, err := uc.minioClient.GetObject(ctx, uc.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}

	objectInfo, err := object.Stat()
	if err != nil {
		// Ensure object is closed on error
		object.Close()
		return nil, minio.ObjectInfo{}, err
	}

	return object, objectInfo, nil
}