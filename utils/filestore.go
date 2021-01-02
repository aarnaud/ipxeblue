package utils

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func NewFileStore(config *Config) *minio.Client {
	ctx := context.Background()
	// Initialize minio client object.
	minioClient, err := minio.New(config.MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioConfig.AccessKey, config.MinioConfig.SecretKey, ""),
		Secure: config.MinioConfig.Secure,
	})
	if err != nil {
		log.Fatalln(err)
	}

	exist, err := minioClient.BucketExists(ctx, config.MinioConfig.BucketName)
	if err != nil {
		log.Fatalln(err)
	}

	if !exist {
		log.Printf("creating bukcet %s \n", config.MinioConfig.BucketName)
		err := minioClient.MakeBucket(ctx, config.MinioConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return minioClient
}

func RemoveRecursive(client *minio.Client, bucketName string, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for objectInfo := range client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    prefix,
		}) {
			if objectInfo.Err != nil {
				log.Println(objectInfo.Err.Error())
			}
			objectsCh <- objectInfo
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for removeObjectError := range client.RemoveObjects(context.Background(), bucketName, objectsCh, opts) {
		return removeObjectError.Err

	}
	return nil
}
