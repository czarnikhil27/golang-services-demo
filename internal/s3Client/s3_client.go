package s3Client

import (
	"bytes"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	svc *s3.S3
}

func NewS3Client(svc *s3.S3) *S3Client {
	return &S3Client{svc: svc}
}
func (c *S3Client) UploadToS3(binaryPath string, projectID string, bucketName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return err
	}
	s3Client := s3.New(sess)

	fileContent, err := os.ReadFile(binaryPath)
	if err != nil {
		return err
	}

	contentType := mime.TypeByExtension(filepath.Ext(binaryPath))

	if _, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fmt.Sprintf("__outputs/%s/app", projectID)),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
	}); err != nil {
		return err
	}
	return nil
}

func (c *S3Client) ListFiles(bucketName, prefix string) ([]string, error) {
	var files []string
	err := c.svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
		return !lastPage
	})
	if err != nil {
		return nil, fmt.Errorf("error listing files in bucket %s: %v", bucketName, err)
	}
	return files, nil
}
