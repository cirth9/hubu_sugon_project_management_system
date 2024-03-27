package min

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strconv"
)

type MinioClient struct {
	c *minio.Client
}

func (c *MinioClient) Put(
	ctx context.Context,
	bucketName string,
	fileName string,
	data []byte,
	size int64,
	contentType string,
) (minio.UploadInfo, error) {
	object, err := c.c.PutObject(
		ctx,
		bucketName,
		fileName,
		bytes.NewBuffer(data),
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return object, err
}

func (c *MinioClient) Compose(
	ctx context.Context,
	bucketName string,
	fileName string,
	totalChunks int,
) (minio.UploadInfo, error) {
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: fileName,
	}
	var srcs []minio.CopySrcOptions
	for i := 1; i <= totalChunks; i++ {
		formatInt := strconv.FormatInt(int64(i), 10)
		src := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: fileName + "_" + formatInt,
		}
		srcs = append(srcs, src)
	}
	object, err := c.c.ComposeObject(
		ctx,
		dst,
		srcs...,
	)
	return object, err
}
func New(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{c: minioClient}, err
}
