package httphandler

import (
	"net/http"
	"time"

	imagestore "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/image_store"
	queueconnector "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/queue_connector"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Maximum size of 10 MB
	r.ParseMultipartForm(10 << 20)

	// Extract file under key image
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Debug().Msgf("Error retrieving file: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()

	userEmail := r.Header.Get(USER_EMAIL_HEADER)

	// This likely means that the auth middleware was not called beforehand
	if len(userEmail) == 0 {
		log.Warn().Msg("Empty user mail header! Something went wrong with the auth middleware")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imageId := uuid.New().String()

	success := imagestore.UploadImage(file, header.Size, imageId)

	if !success {
		log.Warn().Msgf("Could not upload image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	job := queueconnector.Job{
		ImageId:     imageId,
		UserEmail:   userEmail,
		RequestTime: time.Now(),
	}

	ok := queueconnector.AddJobToQueue(job)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
