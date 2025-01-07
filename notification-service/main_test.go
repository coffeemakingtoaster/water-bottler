package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
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
