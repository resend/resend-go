package resend

import (
	"context"
)

type DomainsSvc interface {
	Create(context.Context) error
	Get(context.Context) error
}

type DomainsSvcOp struct {
	client *Client
}

func (s *DomainsSvcOp) Create(ctx context.Context) error {
	path := "domains/"
	println(path)
	println("creating domain")
	return nil
}

func (s *DomainsSvcOp) Get(ctx context.Context) error {
	path := "domains/"
	println(path)
	println("fetching domain info")
	return nil
}
