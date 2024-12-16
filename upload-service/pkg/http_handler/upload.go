package httphandler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	queueconnector "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/queue_connector"
	"github.com/google/uuid"
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
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the file data
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Throw it out (for now)
	// TODO: Upload file to CephFS
	io.Discard.Write(fileBytes)

	imageId := uuid.New().String()
	userEmail := w.Header().Get(USER_EMAIL_HEADER)

	// This likely means that the auth middleware was not called beforehand
	if len(userEmail) == 0 {
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
