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
	AccessKeyID string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
}

func CreateS3DataContainer(dst Destination, name string, creator string, publisher string) error {
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
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}

	mm, err := manifest.NewManifest(creator, publisher)
	if err != nil {
		Log.Error(err)
		return err
	}
	data, _ := json.Marshal(mm)
	minioClient.PutObject(ctx,
		name,
		"META-DATA/manifest.jsonld",
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"})
	return nil
}
