package db

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func UploadToS3(file []byte, fileName, accessKey, secretKey string) error {
	s3Region := os.Getenv("S3_REGION")
	s3Bucket := os.Getenv("S3_BUCKET")
	s3Endpoint := os.Getenv("S3_ENDPOINT")

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(s3Region),
		Endpoint:    aws.String(s3Endpoint),
	})

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(file),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return err
	}

	return nil
}
