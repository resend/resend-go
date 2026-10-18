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

You can view all the examples in the [examples folder](https://github.com/resend/resend-go/tree/main/examples)

## Error handling

For API errors other than HTTP 429, use `errors.As` from the standard library's
`errors` package to inspect the HTTP status, API error name, and message:

```go
_, err := client.Emails.Get("email-id")
var apiError *resend.ResendError
if errors.As(err, &apiError) {
    fmt.Printf("Status: %d, name: %s, message: %s\n",
        apiError.StatusCode, apiError.Name, apiError.Message)
}
```

`StatusCode` comes from the HTTP response. If the error body cannot be read as
JSON, `Name` is empty and `Message` falls back to the HTTP status.
HTTP 429 responses continue to return `*resend.RateLimitError`; use
`errors.Is(err, resend.ErrRateLimit)` to detect them and `errors.As` to access
their rate-limit metadata, including `RetryAfter`.
