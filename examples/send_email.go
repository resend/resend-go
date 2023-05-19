package examples

import (
	"fmt"
	"os"

	"github.com/drish/resend-go"
)

func sendEmailExample() {

	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Send
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

	// Get
	email, err := client.Emails.Get(sent.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", email)

}
