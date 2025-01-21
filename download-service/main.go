package main

import (
	"io"
	"net/http"
	"os"

	"github.com/minio/minio-go"
	log "github.com/rs/zerolog/log"
)

func getHealth(w http.ResponseWriter, r *http.Request, objectStoreConnAvailable bool) {
	log.Info().Msg("Got health check request")
	// Check if the service has access to the object store
	if !objectStoreConnAvailable {
		log.Error().Msg("Object store connection not available")
		http.Error(w, "Not ok", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}

func downloadFile(w http.ResponseWriter, r *http.Request, minioClient *minio.Client, bucketName string) {
	log.Info().Msg("Got download request")
	// Get the file name from the request
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		log.Error().Msg("File name not provided")
		http.Error(w, "file name not provided", http.StatusBadRequest)
		return
	}
	// Get the object
	obj, err := minioClient.GetObject(bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error getting object")
		http.Error(w, "error getting object", http.StatusInternalServerError)
		return
	} else {
		defer obj.Close()
		log.Debug().Str("file", fileName).Msg("Got object")
	}
	// Set the content type
	w.Header().Set("Content-Type", "application/octet-stream")
	// Copy the object to the response writer
	log.Debug().Str("file", fileName).Msg("Downloading file")
	_, err = io.Copy(w, obj)
	if err != nil {
		log.Error().Err(err).Msg("Error copying object")
		http.Error(w, "error copying object", http.StatusInternalServerError)
		return
	}
	log.Debug().Str("file", fileName).Msg("Downloaded file")
}

func main() {
	accessKeyID := os.Getenv("ACCESS_KEY")
	secretAccessKey := os.Getenv("SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")

	objectStoreConnAvailable := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.Error().Err(err).Msg("Error creating minio client")
	}

	// Check if the service has access to the bucket
	_, err = minioClient.BucketExists(bucketName)
	if err == nil {
		objectStoreConnAvailable = true
	} else {
		log.Error().Err(err).Msg("Error checking bucket")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		getHealth(w, r, objectStoreConnAvailable)
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		downloadFile(w, r, minioClient, bucketName)
	})
	log.Info().Msg("Server started")
	err = http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Server encountered error")
}
