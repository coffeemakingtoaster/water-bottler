package main

import (
	"context"
	"math/rand"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minio/minio-go"
	log "github.com/rs/zerolog/log"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(charset[rand.Intn(len(charset))])
	}
	return b.String()
}

func Test_downloadFile(t *testing.T) {
	ctx := context.Background()
	minioContainer, err := tcminio.Run(ctx, "minio/minio:latest", tcminio.WithUsername(generateRandomString(16)), tcminio.WithPassword(generateRandomString(16)))
	if err != nil {
		t.Fatal(err)
	}
	log.Debug().Msg("Minio container started")

	minioEndpoint, err := minioContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	log.Debug().Str("endpoint", minioEndpoint).Msg("Minio endpoint")

	minioClient, err := minio.New(minioEndpoint, minioContainer.Username, minioContainer.Password, false)
	if err != nil {
		log.Error().Err(err).Msg("Error creating minio client")
	}
	log.Debug().Msg("Minio client created")

	// Create a bucket
	bucketName := "test-bucket"
	err = minioClient.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		log.Error().Err(err).Msg("Error creating bucket")
	}
	log.Debug().Str("bucket", bucketName).Msg("Bucket created")

	// Upload a file
	fileName := "test.txt"
	fileContent := "Hello, World!"
	_, err = minioClient.PutObject(bucketName, fileName, strings.NewReader(fileContent), int64(len(fileContent)), minio.PutObjectOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error uploading object")
	}
	log.Debug().Str("file", fileName).Msg("File uploaded")

	// Test the downloadFile function
	uri := strings.Join([]string{"/download?file=", fileName}, "")
	req := httptest.NewRequest("GET", uri, nil)
	rec := httptest.NewRecorder()

	log.Debug().Msg("Testing downloadFile")
	downloadFile(rec, req, minioClient, bucketName)

	if rec.Code != 200 {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	if rec.Body.String() != fileContent {
		t.Errorf("Expected body %s, got %s", fileContent, rec.Body.String())
	}
	log.Debug().Msg("downloadFile test passed")

	// Clean up
	err = minioClient.RemoveObject(bucketName, fileName)
	if err != nil {
		log.Error().Err(err).Msg("Error removing object")
	}

	err = minioClient.RemoveBucket(bucketName)
	if err != nil {
		log.Error().Err(err).Msg("Error removing bucket")
	}

	err = minioContainer.Terminate(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
