# Resend Go SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Build](https://github.com/resend/resend-go/actions/workflows/go.yml/badge.svg)
![Release](https://img.shields.io/github/release/resend/resend-go.svg?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/resend/resend-go/v4.svg)](https://pkg.go.dev/github.com/resend/resend-go/v4)
---

## Installation

To install the Go SDK, simply execute the following command on a terminal:

```
go get github.com/resend/resend-go/v4
```

## Setup

First, you need to get an API key, which is available in the [Resend Dashboard](https://resend.com).

## Example

```go
package main

import (
    "fmt"
    "github.com/resend/resend-go/v4"
)

func main() {
    apiKey := "re_123"

    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        To:      []string{"to@example.com", "you@example.com"},
        From:    "me@example.com",
        Text:    "hello world",
        Subject: "Hello from Golang",
        Cc:      []string{"cc@example.com"},
        Bcc:     []string{"bcc@example.com"},
        ReplyTo: "replyto@example.com",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        panic(err)
    }
    fmt.Println(sent.Id)
}

```

Replace the placeholder API key and email addresses before sending. Use a verified domain for the `From` address.

You can view all the examples in the [examples folder](https://github.com/resend/resend-go/tree/main/examples)
