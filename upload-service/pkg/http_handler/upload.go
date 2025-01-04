package httphandler

import (
	"io"
	"net/http"
	"time"

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
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Debug().Msgf("Error retrieving file: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the file data
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Debug().Msgf("Error reading request file data: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Throw it out (for now)
	// TODO: Upload file to CephFS
	io.Discard.Write(fileBytes)

	imageId := uuid.New().String()
	userEmail := r.Header.Get(USER_EMAIL_HEADER)

	// This likely means that the auth middleware was not called beforehand
	if len(userEmail) == 0 {
		log.Warn().Msg("Empty user mail header! Something went wrong with the auth middleware")
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
