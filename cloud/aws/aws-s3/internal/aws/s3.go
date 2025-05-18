package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitS3Client(cfg aws.Config) {
	S3Client = s3.NewFromConfig(cfg)
}
