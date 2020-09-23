package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/sunet/s3-mm-tool/manifest"
	"encoding/json"
)

var Log = logrus.New()

type Destination struct {
	Endpoint string `json:"endpoint"`,
	Region string `json:"region"`,
	AccessKeyID string `json:"accessKey"`,
	SecretKey string `json:"secretKey"`
}

func CreateS3DataContainer(dst Destination, name string, creator string, publisher string) error {
	ctx := context.Background()

	minioClient, err := minio.New(dst.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(dst.AccessKeyID, dst.SecretKey, ""),
		Secure: true
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
	
	mm := manifest.NewManifest(creator, publisher)
	data, _ := json.Marshal(mm)
	minioClient.PubObject(name, "META-DATA/manifest.jsonld")
}