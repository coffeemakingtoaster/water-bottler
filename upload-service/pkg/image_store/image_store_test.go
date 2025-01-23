package imagestore

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

const TestingBucketName = "testbucket"

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(charset[rand.Intn(len(charset))])
	}
	return b.String()
}

func doesFileExistInMinio(minioEndpoint, minioAccessKeyID, minioSecret, filename string) (bool, error) {
	client, err := minio.New(minioEndpoint, &minio.Options{Secure: false, Creds: credentials.NewStaticV4(minioAccessKeyID, minioSecret, "")})

	if err != nil {
		return false, err
	}
	exists, err := client.BucketExists(context.TODO(), TestingBucketName)

	if err != nil || !exists {
		return false, err
	}

	info, err := client.StatObject(context.TODO(), TestingBucketName, filename, minio.StatObjectOptions{})

	if err != nil || info.Err != nil {
		return false, err
	}

	return true, nil
}

func LoadFile(path string) (*os.File, int64, error) {
	reader, err := os.Open(path)

	if err != nil {
		return nil, 0, err
	}

	stat, err := reader.Stat()

	if err != nil {
		return nil, 0, err
	}

	return reader, stat.Size(), nil
}

func Test_uploadImage(t *testing.T) {
	ctx := context.Background()
	minioContainer, err := tcminio.Run(ctx, "minio/minio:latest", tcminio.WithUsername(generateRandomString(16)), tcminio.WithPassword(generateRandomString(16)))
	if err != nil {
		t.Fatal(err)
	}

	defer minioContainer.Terminate(context.TODO())

	minioEndpoint, err := minioContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	accessKeyID = minioContainer.Username
	secretAccessKey = minioContainer.Password
	endpoint = minioEndpoint
	bucketName = TestingBucketName

	bucketExists = false
	client = nil

	imageId := generateRandomString(32)

	fileReader, size, err := LoadFile("../../testfiles/testimage.jpg")
	if err != nil {
		t.Fatalf("Could not find file due to an error: %s", err.Error())
	}

	success := UploadImage(fileReader, size, imageId)

	if !success {
		t.Fatalf("Did not succeed in uploading image")
	}

	success, err = doesFileExistInMinio(minioEndpoint, minioContainer.Username, minioContainer.Password, imageId)

	if !success || err != nil {
		t.Fatalf("File was not found in minio due to an error: %s", err.Error())
	}
}
