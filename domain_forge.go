package resend

import (
	"context"
	"errors"
	"net/http"
)

// EnableDomainForgeRequest contains params for activating Forge on a domain.
type EnableDomainForgeRequest struct {
	ReturnPath string `json:"return_path,omitempty"`
}

// UpdateDomainForgeRequest contains params for moving the return path of a
// Forge domain that has not verified yet.
type UpdateDomainForgeRequest struct {
	ReturnPath string `json:"return_path"`
}

// VerifyDomainForgeResponse is the slim acknowledgement returned by the
// asynchronous Forge verification, mirroring POST /domains/{id}/verify.
type VerifyDomainForgeResponse struct {
	Object string `json:"object,omitempty"`
	Id     string `json:"id,omitempty"`
}

type DomainForgeSvc interface {
	Enable(domainId string, params *EnableDomainForgeRequest) (Domain, error)
	EnableWithContext(ctx context.Context, domainId string, params *EnableDomainForgeRequest) (Domain, error)
	Update(domainId string, params *UpdateDomainForgeRequest) (Domain, error)
	UpdateWithContext(ctx context.Context, domainId string, params *UpdateDomainForgeRequest) (Domain, error)
	Verify(domainId string) (VerifyDomainForgeResponse, error)
	VerifyWithContext(ctx context.Context, domainId string) (VerifyDomainForgeResponse, error)
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

// UpdateWithContext moves the return path of a Forge domain and re-verifies,
// using the given context. After the domain verifies, the API refuses this
// with a 422 return_path_unavailable error.
func (s *DomainForgeSvcImpl) UpdateWithContext(ctx context.Context, domainId string, params *UpdateDomainForgeRequest) (Domain, error) {
	if params == nil || params.ReturnPath == "" {
		return Domain{}, errors.New("[ERROR]: ReturnPath cannot be empty")
	}

	path := "domains/" + domainId + "/forge"

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return Domain{}, errors.New("[ERROR]: Failed to create DomainForge.Update request")
	}

	resp := new(Domain)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Domain{}, err
	}
	return *resp, nil
}

// Update moves the return path of a Forge domain and re-verifies. After the
// domain verifies, the API refuses this with a 422 return_path_unavailable
// error.
func (s *DomainForgeSvcImpl) Update(domainId string, params *UpdateDomainForgeRequest) (Domain, error) {
	return s.UpdateWithContext(context.Background(), domainId, params)
}

// VerifyWithContext triggers asynchronous Forge verification for a domain,
// using the given context. The response is a slim acknowledgement; poll
// Domains.Get for state.
func (s *DomainForgeSvcImpl) VerifyWithContext(ctx context.Context, domainId string) (VerifyDomainForgeResponse, error) {
	path := "domains/" + domainId + "/forge/verify"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return VerifyDomainForgeResponse{}, errors.New("[ERROR]: Failed to create DomainForge.Verify request")
	}

	resp := new(VerifyDomainForgeResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return VerifyDomainForgeResponse{}, err
	}
	return *resp, nil
}

// Verify triggers asynchronous Forge verification for a domain. The response
// is a slim acknowledgement; poll Domains.Get for state.
func (s *DomainForgeSvcImpl) Verify(domainId string) (VerifyDomainForgeResponse, error) {
	return s.VerifyWithContext(context.Background(), domainId)
}
