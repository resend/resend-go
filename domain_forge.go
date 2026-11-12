package resend

import (
	"context"
	"errors"
	"net/http"
)

// EnableDomainForgeRequest contains params for activating Forge on a domain.
// ReturnPath is optional; when omitted, the server picks a free DNS label.
type EnableDomainForgeRequest struct {
	ReturnPath string `json:"return_path,omitempty"`
}

type DomainForgeSvc interface {
	Enable(domainId string, params *EnableDomainForgeRequest) (Domain, error)
	EnableWithContext(ctx context.Context, domainId string, params *EnableDomainForgeRequest) (Domain, error)
}

type DomainForgeSvcImpl struct {
	client *Client
}

// EnableWithContext activates Forge on a domain and returns the full domain
// object, using the given context. Enabling an already-migrated domain is
// idempotent on the API side and returns its current state.
func (s *DomainForgeSvcImpl) EnableWithContext(ctx context.Context, domainId string, params *EnableDomainForgeRequest) (Domain, error) {
	path := "domains/" + domainId + "/forge"

	if params == nil {
		params = &EnableDomainForgeRequest{}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return Domain{}, errors.New("[ERROR]: Failed to create DomainForge.Enable request")
	}

	resp := new(Domain)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Domain{}, err
	}
	return *resp, nil
}

// Enable activates Forge on a domain and returns the full domain object.
// Enabling an already-migrated domain is idempotent on the API side and
// returns its current state.
func (s *DomainForgeSvcImpl) Enable(domainId string, params *EnableDomainForgeRequest) (Domain, error) {
	return s.EnableWithContext(context.Background(), domainId, params)
}
