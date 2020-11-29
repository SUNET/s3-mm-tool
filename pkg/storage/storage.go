package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/sirupsen/logrus"
	"github.com/sunet/s3-mm-tool/pkg/manifest"
)

var Log = logrus.New()

type Destination struct {
	Endpoint    string `json:"endpoint"`
	Region      string `json:"region"`
	AccessKeyID string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
}

func CreateS3DataContainer(dst Destination, name string, mm manifest.Manifest) error {
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

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:          bytes.NewReader(data),
		ContentLength: int64(len(data)),
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
}

func CreateS3DataContainerMinio(dst Destination, name string, mm manifest.Manifest) error {
	ctx := context.Background()

	minioClient, err := minio.New(dst.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(dst.AccessKeyID, dst.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return err
	}

	err = minioClient.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: dst.Region})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, name)
		if errBucketExists != nil || !exists {
			return err
		}
	}

	data, _ := json.Marshal(mm)
	Log.Debug(string(data))
	info, err := minioClient.PutObject(ctx,
		name,
		".metadata/manifest.jsonld",
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"})

	Log.Debug(info)
	return err
}
