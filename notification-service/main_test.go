package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	queueconnector "github.com/coffeemakingtoaster/water-bottler/notification-service/pkg/queue_connector"
	"github.com/google/uuid"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

type sendMailRequestBody struct {
	Email   string `json:"email"`
	ImageId string `json:"imageid"`
}

func getRequestWithJsonBody(body sendMailRequestBody) *http.Request {
	parsedBody, _ := json.Marshal(&body)
	return httptest.NewRequest("POST", "/send-mail", bytes.NewReader(parsedBody))
}

func setupTemplate() {
	MAIL_TEMPLATE, _ = template.New("Test Template").Parse("From: {{.From}}\r\n" +
		"To: {{.Email}}\r\n" +
		"Subject: Testing!\r\n" +
		"\r\n" +
		"{{.ImageId}}\r\n")
}

func setupMockSmtpServer() *smtpmock.Server {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := server.Start(); err != nil {
		panic(err)
	}
	return server
}

func gracefullyShutdownSMTPSever(server *smtpmock.Server) {
	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}
}

func Test_successfulSmtpSend(t *testing.T) {
	s := setupMockSmtpServer()
	defer gracefullyShutdownSMTPSever(s)

	setupTemplate()
	w := httptest.NewRecorder()

	expectedMail := fmt.Sprintf("%s@testing.test", uuid.New().String())
	expectedImageId := uuid.New().String()

	SMTP_SERVER_URL = fmt.Sprintf("localhost:%d", s.PortNumber())

	sendMail(w, getRequestWithJsonBody(sendMailRequestBody{Email: expectedMail, ImageId: expectedImageId}))

	msgs := s.MessagesAndPurge()

	foundMail := false
	foundImageId := false

	for _, msg := range msgs {
		foundMail = strings.Contains(msg.MsgRequest(), fmt.Sprintf("To: %s", expectedMail)) || foundMail
		foundImageId = strings.Contains(msg.MsgRequest(), expectedImageId) || foundImageId
	}

	if !foundMail {
		t.Error("Recipient mail not found in email headers")
	}

	if !foundImageId {
		t.Error("ImageId not found in mail content")
	}
}

func Test_smtpServerNotReachable(t *testing.T) {
	setupTemplate()

	w := httptest.NewRecorder()

	SMTP_SERVER_URL = fmt.Sprintf("localhost:1234")

	sendMail(w, getRequestWithJsonBody(sendMailRequestBody{Email: "test@test.test", ImageId: "testing"}))

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected statuscode %d but received %d", http.StatusInternalServerError, w.Result().StatusCode)
	}
}

func Test_sendMailNoBody(t *testing.T) {
	setupTemplate()
	w := httptest.NewRecorder()

	SMTP_SERVER_URL = fmt.Sprintf("localhost:1234")

	sendMail(w, httptest.NewRequest("POST", "/send-mail", nil))

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d but received %d", http.StatusBadRequest, w.Result().StatusCode)
	}
}

func Test_finishedJobConsumer(t *testing.T) {
	ctx := context.Background()
	rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3-management-alpine")
	if err != nil {
		t.Errorf("Could not start container due to an error %s", err.Error())
	}

	defer func() {
		if err := rabbitContainer.Terminate(context.TODO()); err != nil {
			t.Errorf("failed to terminate container: %s", err)
		}
	}()

	queueconnector.QueueUrl, err = rabbitContainer.AmqpURL(ctx)

	if err != nil {
		t.Errorf("Could not retrieve connection uri to rabbitmq due to an error: %s", err.Error())
	}

	testJob := queueconnector.FinishedJob{
		ImageId:   "test.png",
		UserEmail: "test@water-bottler.com",
	}

	publishedJob := queueconnector.AddJobToQueue(testJob)
	if !publishedJob {
		t.Error("Could not publish job to queue")
	} else {
		t.Log("Test job published")
	}

	finishedJobs, success := queueconnector.ConsumeJobFromQueue()
	if !success {
		t.Error("Could not consume jobs from queue")
	}

	var wrongJobsCount int = 0
	for job := range finishedJobs {
		if job.UserEmail != testJob.UserEmail || job.ImageId != testJob.ImageId {
			t.Errorf("Received job did not match expected job. Expected: %v\tGot: %v", testJob, job)
			wrongJobsCount++
		} else {
			t.Log("Received test job")
			break
		}

		if wrongJobsCount > 5 {
			t.Error("Received more than one wrong job, aborting test \n Maybe the queue is used by other applications?")
			break
		}
	}
}
