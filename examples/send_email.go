package main

import (
	"fmt"
	"os"

	"github.com/drish/resend-go"
)

// Rename to main
func sendEmail() {

	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		panic("Api Key is missing")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		To:      []string{"carlosderich@gmail.com", "derich@thinkdataworks.com"},
		From:    "r@recomendo.io",
		Text:    "hello world",
		Subject: "Hello from Golang",
		Cc:      []string{"d.erich@hotmail.com"},
		Bcc:     []string{"d.erich@hotmail.com"},
		ReplyTo: "carlosderich@gmail.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
