package mails

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"os"
	"path/filepath"

	"github.com/YahiaJouini/chat-app-backend/internal/config"
	"gopkg.in/gomail.v2"
)

// GenerateVerificationCode Generate a random 6-digit verification code
func GenerateVerificationCode() (string, error) {
	maximum := 999999
	minimum := 100000
	n, err := rand.Int(rand.Reader, big.NewInt(int64(maximum-minimum+1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()+int64(minimum)), nil
}

type Result struct {
	Err error
}

func Success() Result {
	return Result{}
}

func Failure(err error) Result {
	return Result{Err: err}
}

func SendMail(sendTo string, code string) Result {
	// get html
	var body bytes.Buffer
	currentDir, _ := os.Getwd()
	templatePath := filepath.Join(currentDir, "pkg", "email", "verification.html")

	tmpl, err := template.New("verification.html").ParseFiles(templatePath)
	if err != nil {
		fmt.Println("failed to parse template: ", err)
		return Failure(err)
	}

	err = tmpl.Execute(&body, struct{ Code string }{Code: code})
	if err != nil {
		fmt.Println("error executing template", err)
		return Failure(err)
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	emailSender, _ := config.GetEnv("EMAIL_SENDER")
	emailPassword, _ := config.GetEnv("EMAIL_PASSWORD")

	message := gomail.NewMessage()
	message.SetHeader("From", emailSender)
	message.SetHeader("To", sendTo)
	message.SetHeader("Subject", fmt.Sprintf("%v is your verification code", code))
	message.SetBody("text/html", body.String())

	// send email
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, emailPassword)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("error sending mail", err)
		return Failure(err)
	}
	return Success()
}
