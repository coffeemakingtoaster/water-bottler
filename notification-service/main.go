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

	"github.com/rs/zerolog/log"
)

var SMTP_SERVER_URL string

var MAIL_TEMPLATE *template.Template

type MailRequestData struct {
	Email   string
	ImageId string
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got health request")
	io.WriteString(w, "ok")
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var requestData MailRequestData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	from := "notification@water-bottler.com"
	password := "test123"

	var message bytes.Buffer

	MAIL_TEMPLATE.Execute(&message, requestData)

	auth := smtp.PlainAuth("", from, password, SMTP_SERVER_URL)

	err = smtp.SendMail(SMTP_SERVER_URL, auth, from, []string{requestData.Email}, message.Bytes())
	if err != nil {
		log.Warn().Msgf("Could not send mail due to an error: %s", err.Error())
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

	var err error
	MAIL_TEMPLATE, err = template.ParseFiles("./mail.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Could not parse mail template due to an error: %s", err.Error()))
	}

	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/send-mail", getHealth)
	addr := fmt.Sprintf("%s:%d", interfaceIP, interfacePort)
	log.Info().Msgf("Starting authentication service on %s", addr)
	err = http.ListenAndServe(addr, nil)
	fmt.Printf("Server encountered error: %v", err)
}
