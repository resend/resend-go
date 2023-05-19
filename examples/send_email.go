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
		To:      []string{"to@example.com", "to2@example.com"},
		From:    "from@example.com",
		Text:    "hello world",
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"ccc@example.com"},
		ReplyTo: "to@example.com",
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
