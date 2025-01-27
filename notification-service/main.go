package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	queueconnector "github.com/coffeemakingtoaster/water-bottler/notification-service/pkg/queue_connector"
	"github.com/rs/zerolog/log"
)

var SMTP_SERVER_URL string

var SMTP_SERVER_USERNAME string
var SMTP_SERVER_PASSWORD string

var MAIL_TEMPLATE *template.Template

type MailTemplateData struct {
	Data      queueconnector.FinishedJob
	From      string
	SourceUri string
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Got health request")
	io.WriteString(w, "ok")
}

func sendMail(job queueconnector.FinishedJob) bool {

	if len(job.UserEmail) == 0 || len(job.ImageId) == 0 {
		log.Warn().Msg("Received faulty job entity")
		return false
	}

	templateData := MailTemplateData{
		Data:      job,
		From:      "notification@water-bottler.com",
		SourceUri: os.Getenv("SOURCE_URI"),
	}

	var message bytes.Buffer

	MAIL_TEMPLATE.Execute(&message, templateData)

	var err error

	if SMTP_SERVER_PASSWORD != "" && SMTP_SERVER_USERNAME != "" {
		// For basicauth the smtp server url cannot contain port
		cleanURL := strings.Split(SMTP_SERVER_URL, ":")
		auth := smtp.PlainAuth("water-bottler-mail", SMTP_SERVER_USERNAME, SMTP_SERVER_PASSWORD, cleanURL[0])
		err = smtp.SendMail(SMTP_SERVER_URL, auth, templateData.From, []string{templateData.Data.UserEmail}, message.Bytes())
	} else {
		err = smtp.SendMail(SMTP_SERVER_URL, nil, templateData.From, []string{templateData.Data.UserEmail}, message.Bytes())
	}

	if err != nil {
		log.Warn().Msgf("Could not send mail due to an error: %s", err.Error())
		return false
	}
	log.Debug().Msg("Email Sent Successfully!")
	return true
}

func jobConsumer(queueConnectionString string) {
	log.Debug().Msgf("Starting job consumer")
	queueconnector.QueueUrl = queueConnectionString
	if len(queueconnector.QueueUrl) == 0 {
		log.Fatal().Msg("No queue url specified via the QUEUE_URL env variable!")
	}
	log.Debug().Msgf("Queue url is %s", queueconnector.QueueUrl)

	finishedJobs, success := queueconnector.ConsumeJobFromQueue()
	if !success {
		log.Fatal().Msg("Could not consume jobs from queue")
	}

	for job := range finishedJobs {
		log.Debug().Msgf("Received job: %v", job)
		sendMail(job)
	}
}

func main() {
	interfaceIP := "0.0.0.0"
	interfacePort := 8080

	SMTP_SERVER_URL = os.Getenv("SMTP_SERVER_URL")
	if len(SMTP_SERVER_URL) == 0 {
		panic("No smtp server specified via the 'SMTP_SERVER_URL' env variable!")
	}

	// No check if they are actually set because SMTP servers may have no auth required
	SMTP_SERVER_USERNAME = os.Getenv("SMTP_SERVER_USERNAME")
	SMTP_SERVER_PASSWORD = os.Getenv("SMTP_SERVER_PASSWORD")

	log.Debug().Msgf("SMTP server configured to be %s", SMTP_SERVER_URL)

	var err error
	MAIL_TEMPLATE, err = template.ParseFiles("./mail.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Could not parse mail template due to an error: %s", err.Error()))
	}

	// Start the job consumer and send mail notification if a job is finished
	go jobConsumer(os.Getenv("QUEUE_URL"))

	http.HandleFunc("/health", getHealth)
	addr := fmt.Sprintf("%s:%d", interfaceIP, interfacePort)
	log.Info().Msgf("Starting notification service on %s", addr)
	err = http.ListenAndServe(addr, nil)
	log.Fatal().Msgf("Server encountered error: %v", err)
}
