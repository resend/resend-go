package examples

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drish/resend-go"
)

func sendWithAttachments() {

	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		panic("Api Key is missing")
	}

	// Read attachment file
	pwd, _ := os.Getwd()
	f, err := ioutil.ReadFile(pwd + "/resources/invoice.pdf")
	if err != nil {
		panic(err)
	}

	client := resend.NewClient(apiKey)

	// Create attachment object
	pdfAttachment := &resend.Attachment{
		Content:  string(f),
		Filename: "invoice.pdf",
	}

	params := &resend.SendEmailRequest{
		To:          []string{"carlosderich@gmail.com", "derich@thinkdataworks.com"},
		From:        "r@recomendo.io",
		Text:        "take a look at the file I just sent you",
		Subject:     "Email with attachment",
		Attachments: []resend.Attachment{*pdfAttachment},
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
