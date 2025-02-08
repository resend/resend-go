# Resend Go SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Build](https://github.com/resend/resend-go/actions/workflows/go.yml/badge.svg)
![Release](https://img.shields.io/github/release/resend/resend-go.svg?style=flat-square)
---

## Installation

To install the Go SDK, simply execute the following command on a terminal:

```
go get github.com/resend/resend-go/v2
```

## Setup

First, you need to get an API key, which is available in the [Resend Dashboard](https://resend.com).

## Example

```go
import (
    "fmt"
    "github.com/resend/resend-go/v2"
)

func main() {
    apiKey := "re_123"

    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        To:      []string{"to@example", "you@example.com"},
        From:    "me@exemple.io",
        Text:    "hello world",
        Subject: "Hello from Golang",
        Cc:      []string{"cc@example.com"},
        Bcc:     []string{"cc@example.com"},
        ReplyTo: "replyto@example.com",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        panic(err)
    }
    fmt.Println(sent.Id)
}

```
```go
import (
    "fmt"
    "github.com/resend/resend-go/v2"
)

func main() {
    apiKey := "re_123"

    client := resend.NewClient(apiKey)
    // No need to create objects manually
    sent, err := client.Emails.SendWithOption(
        WithSendEmailTo([]string{"to@example", "you@example.com"}),
        WithSendEmailFrom("me@exemple.io"),
        WithSendEmailText("hello world"),
        WithSendEmailSubject("Hello from Golang"),
        WithSendEmailCc([]string{"cc@example.com"}),
        WithSendEmailBcc([]string{"cc@example.com"}),
        WithSendEmailReplyTo("replyto@example.com"),
        WithSendEmailScheduledAt(""),
        WithSendEmailAttachment(
            WithAttachmentPath("/test"),
            WithAttachmentContent([]byte{}),
            WithAttachmentFilename("test.txt"),
        ),
        WithSendEmailAttachment(
            WithAttachmentPath("/test"),
            WithAttachmentContent([]byte{}),
            WithAttachmentFilename("test2.txt"),
        ),
        )
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(sent.Id)
}
```

You can view all the examples in the [examples folder](https://github.com/resend/resend-go/tree/main/examples)
