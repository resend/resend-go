package resend

import (
	"encoding/json"
	"errors"
	"net/http"
)

type DomainsSvc interface {
	Create(*CreateDomainRequest) (CreateDomainResponse, error)
	Verify(domainId string) (bool, error)
	List() (ListDomainsResponse, error)
	Get(domainId string) (Domain, error)
	Remove(domainId string) (bool, error)
}

type DomainsSvcImpl struct {
	client *Client
}

type CreateDomainRequest struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
}

type CreateDomainResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	CreatedAt   string   `json:"createdAt"`
	Status      string   `json:"status"`
	Records     []Record `json:"records"`
	Region      string   `json:"region"`
	DnsProvider string   `json:"dnsProvider"`
}

type ListDomainsResponse struct {
	Data []Domain `json:"data"`
}

type Domain struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Region    string `json:"region"`
}

type Record struct {
	Record   string      `json:"record"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Ttl      string      `json:"ttl"`
	Status   string      `json:"status"`
	Value    string      `json:"value"`
	Priority json.Number `json:"priority,omitempty"`
}

// Create creates a new Domain entry based on the given params
// https://resend.com/docs/api-reference/domains/create-domain
func (s *DomainsSvcImpl) Create(params *CreateDomainRequest) (CreateDomainResponse, error) {
	path := "domains"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, params)
	if err != nil {
		return CreateDomainResponse{}, errors.New("[ERROR]: Failed to create Domains.Create request")
	}

	// Build response recipient obj
	domainsResp := new(CreateDomainResponse)

	// Send Request
	_, err = s.client.Perform(req, domainsResp)

	if err != nil {
		return CreateDomainResponse{}, err
	}

	return *domainsResp, nil
}

// Verify verifies a given domain Id
// https://resend.com/docs/api-reference/domains/verify-domain
func (s *DomainsSvcImpl) Verify(domainId string) (bool, error) {
	path := "domains/" + domainId + "/verify"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return false, errors.New("[ERROR]: Failed to create Domains.Verify request")
	}

	// Send Request
	_, err = s.client.Perform(req, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}

// List returns the list of all doamins
// https://resend.com/docs/api-reference/domains/list-domains
func (s *DomainsSvcImpl) List() (ListDomainsResponse, error) {
	path := "domains"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return ListDomainsResponse{}, errors.New("[ERROR]: Failed to create Domains.Verify request")
	}

	domains := new(ListDomainsResponse)

	// Send Request
	_, err = s.client.Perform(req, domains)

	if err != nil {
		return ListDomainsResponse{}, err
	}

	return *domains, nil
}

// Remove removes a given domain entry by id
// https://resend.com/docs/api-reference/domains/delete-domain
func (s *DomainsSvcImpl) Remove(domainId string) (bool, error) {
	path := "domains/" + domainId

	// Prepare request
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return false, errors.New("[ERROR]: Failed to create Domains.Remove request")
	}

	// Send Request
	_, err = s.client.Perform(req, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Get retrieves a domain object
// https://resend.com/docs/api-reference/domains/get-domain
func (s *DomainsSvcImpl) Get(domainId string) (Domain, error) {
	path := "domains/" + domainId

	// Prepare request
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return Domain{}, errors.New("[ERROR]: Failed to create Domains.Get request")
	}

	domain := new(Domain)

	// Send Request
	_, err = s.client.Perform(req, domain)

	if err != nil {
		return Domain{}, err
	}

	return *domain, nil
}
