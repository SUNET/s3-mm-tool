package storage

import (
	"bytes"
	"context"
	"encoding/json"

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
