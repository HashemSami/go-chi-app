package main

import (
	"fmt"

	"github.com/go-mail/mail"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "user"
	password = "pass"
)

func main() {
	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)

	// msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)

	sender, err := dialer.Dial()
	if err != nil {
		panic(err)
	}

	defer sender.Close()
	sender.Send(from, []string{to}, msg)
	fmt.Println("msg sent>>>")
}
