package repository

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

func (r *Repository) UploadFile(
	reader io.Reader,
	size int64,
	filename string,
	contentType string,
) error {
	_, err := r.minio.PutObject(
		context.Background(),
		r.bucket,
		filename,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	return err
}
