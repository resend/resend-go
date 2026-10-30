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

Use `errors.As` from the standard library's `errors` package to inspect typed
API errors. HTTP 429 responses return `*resend.RateLimitError`; other HTTP error
responses return `*resend.ResendError`. Handle other failures, such as network
errors or canceled requests, separately:

```go
_, err := client.Emails.Get("email-id")
if err != nil {
    var rateLimit *resend.RateLimitError
    var apiError *resend.ResendError

    switch {
    case errors.As(err, &rateLimit):
        fmt.Printf("Rate limited: %s (Retry-After: %s)\n",
            rateLimit.Message, rateLimit.RetryAfter)
    case errors.As(err, &apiError):
        fmt.Printf("Status: %d, name: %s, message: %s\n",
            apiError.StatusCode, apiError.Name, apiError.Message)
    default:
        fmt.Printf("Request failed: %v\n", err)
    }
}
```

`ResendError.StatusCode` comes from the HTTP response. If the error body cannot
be read as JSON, `Name` is empty and `Message` falls back to the HTTP status.
`RateLimitError` includes rate-limit metadata such as `RetryAfter` when the
corresponding response headers are present. You can also detect rate limits
with `errors.Is(err, resend.ErrRateLimit)`.

The SDK does not automatically retry requests. The application decides whether
and when to retry.
