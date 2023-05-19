# Resend Go SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Build](https://github.com/drish/resend-go/actions/workflows/go.yml/badge.svg)
---

## Installation

To install the Go SDK, simply execute the following command on a terminal:

```
go get https://github.com/drish/resend-go
```

## Setup

First, you need to get an API key, which is available in the [Resend Dashboard](https://resend.com).

## Example

```go
import "github.com/drish/resend-go"

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
    fmt.Errorf(err)
}
fmt.Println(sent.Id)
```

You can view all the examples in the [examples folder](https://github.com/drish/resend-go/tree/main/examples)