package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

var SMTP_SERVER_URL string

var SMTP_SERVER_USERNAME string
var SMTP_SERVER_PASSWORD string

var MAIL_TEMPLATE *template.Template

type MailRequestData struct {
	Email   string `json:"email"`
	ImageId string `json:"imageid"`
	From    string
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got health request")
	io.WriteString(w, "ok")
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var requestData MailRequestData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Debug().Msgf("Could not decode body due to an error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData.From = "notification@water-bottler.com"

	var message bytes.Buffer

	MAIL_TEMPLATE.Execute(&message, requestData)

	if SMTP_SERVER_PASSWORD != "" && SMTP_SERVER_USERNAME != "" {
		// For basicauth the smtp server url cannot contain port
		cleanURL := strings.Split(SMTP_SERVER_URL, ":")
		auth := smtp.PlainAuth("water-bottler-mail", SMTP_SERVER_USERNAME, SMTP_SERVER_PASSWORD, cleanURL[0])
		err = smtp.SendMail(SMTP_SERVER_URL, auth, requestData.From, []string{requestData.Email}, message.Bytes())
	} else {
		err = smtp.SendMail(SMTP_SERVER_URL, nil, requestData.From, []string{requestData.Email}, message.Bytes())
	}

	if err != nil {
		log.Warn().Msgf("Could not send mail due to an error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("Email Sent Successfully!")
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

	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/send-mail", sendMail)
	addr := fmt.Sprintf("%s:%d", interfaceIP, interfacePort)
	log.Info().Msgf("Starting notification service on %s", addr)
	err = http.ListenAndServe(addr, nil)
	fmt.Printf("Server encountered error: %v", err)
}
