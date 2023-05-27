package examples

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/resendlabs/resend-go"
)

func withAttachments() {

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

	// Create attachments objects
	pdfAttachmentFromLocalFile := &resend.Attachment{
		Content:  string(f),
		Filename: "invoice1.pdf",
	}

	pdfAttachmentFromRemotePath := &resend.Attachment{
		Path:     "https://github.com/resendlabs/resend-go/raw/main/resources/invoice.pdf",
		Filename: "invoice2.pdf",
	}

	params := &resend.SendEmailRequest{
		To:          []string{"carlosderich@gmail.com", "d.erich@hotmail.com"},
		From:        "from@recomendo.io",
		Text:        "email with attachments !!",
		Html:        "<strong>email with attachments !!</strong>",
		Subject:     "Email with attachment",
		Attachments: []resend.Attachment{*pdfAttachmentFromLocalFile, *pdfAttachmentFromRemotePath},
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
