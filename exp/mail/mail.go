package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/HashemSami/go-chi-app/models"
	"github.com/joho/godotenv"
)

func main() {
	// load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// get
	conf := models.SMTPConfig{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
	}

	es := models.NewEmailService(conf)

	// err := es.Send(email)
	err = es.ForgotPassword("hash@hash.com",
		"https://lenslocked.com/reset-pw?token=abc123",
	)
	if err != nil {
		panic(err)
	}

	// email := models.Email{
	// 	From:      "test@lenslocked.com",
	// 	To:        "jon@calhoun.io",
	// 	Subject:   "This is a test email",
	// 	PlainText: "This is the body of the email",
	// 	HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	// }

	// msg := mail.NewMessage()
	// msg.SetHeader("From", from)
	// msg.SetHeader("To", to)
	// msg.SetHeader("Subject", subject)
	// msg.SetBody("text/plain", plaintext)
	// msg.AddAlternative("text/html", html)

	// // msg.WriteTo(os.Stdout)

	// dialer := mail.NewDialer(host, port, username, password)

	// sender, err := dialer.Dial()
	// if err != nil {
	// 	panic(err)
	// }

	// defer sender.Close()
	// // used for multiple messages send
	// sender.Send(from, []string{to}, msg)
	// sender.Send(from, []string{to}, msg)
	// sender.Send(from, []string{to}, msg)
	// sender.Send(from, []string{to}, msg)
	fmt.Println("msg sent>>>")
}
