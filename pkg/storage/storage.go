package storage

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/sirupsen/logrus"
	"github.com/sunet/s3-mm-tool/pkg/manifest"
)

var Log = logrus.New()

// Destination example
type Destination struct {
	Endpoint    string `json:"endpoint" example:"https://s3.example.com/"`
	Region      string `json:"region" example:"example-region"`
	AccessKeyID string `json:"access_key" example:"access1"`
	SecretKey   string `json:"secret_key" example:"secret1"`
}

func CreateS3DataContainer(dst Destination, name string, mm manifest.ManifestInfo) error {
	bucket := aws.String(name)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(dst.AccessKeyID, dst.SecretKey, ""),
		Endpoint:         aws.String(dst.Endpoint),
		Region:           aws.String(dst.Region),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket, // Required
	}

	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(mm)
	dataLen := int64(len(data))
	Log.Debug(string(data))
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:          bytes.NewReader(data),
		ContentLength: &dataLen,
		ContentType:   aws.String("application/json"),
		Bucket:        bucket,
		Key:           aws.String(".metadata/manifest.jsonld"),
	})
	if err != nil {
		return err
	}

	_, err = s3Client.PutBucketAcl(&s3.PutBucketAclInput{
		Bucket:           bucket,
		GrantFullControl: aws.String(os.Getenv("S3_API_USER")),
	})
	if err != nil {
		return err
	}

	return nil
}
