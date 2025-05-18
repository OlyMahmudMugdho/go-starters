package service

import (
	"aws-s3/internal/aws"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadFile(ctx context.Context, bucketName, key string, file multipart.File) error {
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	_, err = aws.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func DownloadFile(ctx context.Context, bucketName, key string) (io.ReadCloser, error) {
	output, err := aws.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return output.Body, nil
}
func DeleteFile(ctx context.Context, bucketName, key string) error {
	_, err := aws.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
