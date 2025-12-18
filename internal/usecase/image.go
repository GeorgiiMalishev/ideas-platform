package usecase

import (
	"context"
	"mime/multipart"
	"github.com/minio/minio-go/v7"
)

type ImageUsecase interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error)
	CreateBucket(ctx context.Context) error
	GetImage(ctx context.Context, objectName string) (*minio.Object, minio.ObjectInfo, error)
}
