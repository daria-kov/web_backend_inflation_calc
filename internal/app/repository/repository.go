package repository

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	minio  *minio.Client
	bucket string
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("MINIO_ACCESS_KEY"),
				os.Getenv("MINIO_SECRET_KEY"),
				"",
			),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}

	bucket := os.Getenv("MINIO_BUCKET")

	exists, err := minioClient.BucketExists(
		context.Background(),
		bucket,
	)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(
			context.Background(),
			bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return nil, err
		}
	}

	return &Repository{
		db:     db,
		minio:  minioClient,
		bucket: bucket,
	}, nil
}
