package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"aws-s3/internal/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func CreateBucket(ctx context.Context, name string, region string) error {
	_, err := aws.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &name,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) || errors.As(err, &exists) {
			log.Printf("Bucket %s already exists or owned.\n", name)
			return nil
		}
		return err
	}

	waiter := s3.NewBucketExistsWaiter(aws.S3Client)
	err = waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: &name}, time.Minute)
	if err != nil {
		log.Printf("Waiter failed for bucket %s\n", name)
		return err
	}

	return nil
}

func ListBuckets(ctx context.Context) ([]string, error) {
	output, err := aws.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "AccessDenied" {
			return nil, errors.New("access denied")
		}
		return nil, err
	}

	var bucketNames []string
	for _, b := range output.Buckets {
		bucketNames = append(bucketNames, *b.Name)
	}
	return bucketNames, nil
}

func DeleteBucket(ctx context.Context, bucket string) error {
	_, err := aws.S3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket: %w", err)
	}

	waiter := s3.NewBucketNotExistsWaiter(aws.S3Client)
	err = waiter.Wait(ctx, &s3.HeadBucketInput{
		Bucket: &bucket,
	}, 2*time.Minute)

	if err != nil {
		return fmt.Errorf("error waiting for bucket to be deleted: %w", err)
	}
	return nil
}

func UpdateBucket(ctx context.Context, oldName, newName, newRegion string) error {
	if err := DeleteBucket(ctx, oldName); err != nil {
		return fmt.Errorf("failed to delete old bucket: %w", err)
	}
	if err := CreateBucket(ctx, newName, newRegion); err != nil {
		return fmt.Errorf("failed to create new bucket: %w", err)
	}
	return nil
}
