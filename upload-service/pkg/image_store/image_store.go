package imagestore

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

var accessKeyID string
var secretAccessKey string
var endpoint string
var bucketName string

var client *minio.Client

var objectStoreAvailable = false

var bucketExists = false

func init() {
	accessKeyID = os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey = os.Getenv("MINIO_SECRET_KEY")
	endpoint = os.Getenv("MINIO_ENDPOINT")
	bucketName = os.Getenv("MINIO_BUCKET_NAME")
}

func getClient() *minio.Client {
	if client == nil {
		var err error
		client, err = minio.New(endpoint, &minio.Options{Secure: false, Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, "")})
		if err != nil {
			log.Warn().Msgf("Could not create minio client due to an error: %s", err.Error())
			objectStoreAvailable = false
			return nil
		}
	}
	if bucketExists {
		objectStoreAvailable = true
		return client
	}

	exists, err := client.BucketExists(context.TODO(), bucketName)

	if err != nil {
		log.Warn().Msgf("Could check bucket due to an error: %s", err.Error())
		objectStoreAvailable = false
		return nil
	}

	if exists {
		objectStoreAvailable = true
		return client
	}

	err = client.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{})

	if err != nil {
		log.Warn().Msgf("Could not create bucket due to an error: %s", err.Error())
		objectStoreAvailable = false
		return nil
	}
	objectStoreAvailable = true
	return client
}

func UploadImage(file io.Reader, size int64, imageid string) bool {
	c := getClient()
	if c == nil {
		return false
	}
	info, err := c.PutObject(context.TODO(), bucketName, imageid, file, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Warn().Msgf("Could not upload image due to an error: %s", err.Error())
		objectStoreAvailable = false
		return false
	}
	log.Debug().Msgf("Uploaded image of size %d", info.Size)
	return true
}

func IsHealthy() bool {
	return objectStoreAvailable
}
