package resend

import (
	"context"
	"encoding/json"
	"net/http"
)

type TlsOption = string //nolint:revive

const (
	Enforced      TlsOption = "enforced"
	Opportunistic TlsOption = "opportunistic"
)

type DomainsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateDomainRequest) (CreateDomainResponse, error)
	Create(params *CreateDomainRequest) (CreateDomainResponse, error)
	VerifyWithContext(ctx context.Context, domainId string) (bool, error) //nolint:revive
	Verify(domainId string) (bool, error)                                 //nolint:revive
	ListWithOptions(ctx context.Context, options *ListOptions) (ListDomainsResponse, error)
	ListWithContext(ctx context.Context) (ListDomainsResponse, error)
	List() (ListDomainsResponse, error)
	GetWithContext(ctx context.Context, domainId string) (Domain, error)                                 //nolint:revive
	Get(domainId string) (Domain, error)                                                                 //nolint:revive
	RemoveWithContext(ctx context.Context, domainId string) (bool, error)                                //nolint:revive
	Remove(domainId string) (bool, error)                                                                //nolint:revive
	UpdateWithContext(ctx context.Context, domainId string, params *UpdateDomainRequest) (Domain, error) //nolint:revive
	Update(domainId string, params *UpdateDomainRequest) (Domain, error)                                 //nolint:revive
}

type DomainsSvcImpl struct {
	client *Client
}

type CreateDomainRequest struct {
	Name             string `json:"name"`
	Region           string `json:"region,omitempty"`
	CustomReturnPath string `json:"custom_return_path,omitempty"`
}

type CreateDomainResponse struct {
	Id          string   `json:"id"` //nolint:revive
	Name        string   `json:"name"`
	CreatedAt   string   `json:"createdAt"`
	Status      string   `json:"status"`
	Records     []Record `json:"records"`
	Region      string   `json:"region"`
	DnsProvider string   `json:"dnsProvider"` //nolint:revive
}

type ListDomainsResponse struct {
	Object  string   `json:"object"`
	Data    []Domain `json:"data"`
	HasMore bool     `json:"has_more"`
}

type UpdateDomainRequest struct {
	OpenTracking  bool      `json:"open_tracking,omitempty"`
	ClickTracking bool      `json:"click_tracking,omitempty"`
	Tls           TlsOption `json:"tls,omitempty"` //nolint:revive
}

type Domain struct {
	Id        string   `json:"id,omitempty"` //nolint:revive
	Object    string   `json:"object,omitempty"`
	Name      string   `json:"name,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	Status    string   `json:"status,omitempty"`
	Region    string   `json:"region,omitempty"`
	Records   []Record `json:"records,omitempty"`
}

type Record struct {
	Record   string      `json:"record"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Ttl      string      `json:"ttl"` //nolint:revive
	Status   string      `json:"status"`
	Value    string      `json:"value"`
	Priority json.Number `json:"priority,omitempty"`
}

// UpdateWithContext updates an existing Domain entry based on the given params
// https://resend.com/docs/api-reference/domains/update-domain
func (s *DomainsSvcImpl) UpdateWithContext(ctx context.Context, domainId string, params *UpdateDomainRequest) (Domain, error) { //nolint:revive
	path := "domains/" + domainId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return Domain{}, ErrFailedToCreateDomainsUpdateRequest
	}

	domainUpdatedResp := new(Domain)

	// Send Request
	_, err = s.client.Perform(req, domainUpdatedResp) //nolint:bodyclose
	if err != nil {
		return Domain{}, err
	}

	return *domainUpdatedResp, nil
}

// Update is a wrapper around UpdateWithContext
func (s *DomainsSvcImpl) Update(domainId string, params *UpdateDomainRequest) (Domain, error) { //nolint:revive
	return s.UpdateWithContext(context.Background(), domainId, params)
}

// CreateWithContext creates a new Domain entry based on the given params
// https://resend.com/docs/api-reference/domains/create-domain
func (s *DomainsSvcImpl) CreateWithContext(ctx context.Context, params *CreateDomainRequest) (CreateDomainResponse, error) {
	path := "domains"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateDomainResponse{}, ErrFailedToCreateDomainsCreateRequest
	}

	// Build response recipient obj
	domainsResp := new(CreateDomainResponse)

	// Send Request
	_, err = s.client.Perform(req, domainsResp) //nolint:bodyclose
	if err != nil {
		return CreateDomainResponse{}, err
	}

	return *domainsResp, nil
}

// Create creates a new Domain entry based on the given params
// https://resend.com/docs/api-reference/domains/create-domain
func (s *DomainsSvcImpl) Create(params *CreateDomainRequest) (CreateDomainResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// VerifyWithContext verifies a given domain Id
// https://resend.com/docs/api-reference/domains/verify-domain
func (s *DomainsSvcImpl) VerifyWithContext(ctx context.Context, domainId string) (bool, error) { //nolint:revive
	path := "domains/" + domainId + "/verify"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return false, ErrFailedToCreateDomainsVerifyRequest
	}

	// Send Request
	_, err = s.client.Perform(req, nil) //nolint:bodyclose
	if err != nil {
		return false, err
	}

	return true, nil
}

// Verify verifies a given domain Id
// https://resend.com/docs/api-reference/domains/verify-domain
func (s *DomainsSvcImpl) Verify(domainId string) (bool, error) { //nolint:revive
	return s.VerifyWithContext(context.Background(), domainId)
}

// ListWithOptions returns the list of all domains with pagination options
// https://resend.com/docs/api-reference/domains/list-domains
func (s *DomainsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListDomainsResponse, error) {
	path := "domains" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListDomainsResponse{}, ErrFailedToCreateDomainsListRequest
	}

	domains := new(ListDomainsResponse)

	// Send Request
	_, err = s.client.Perform(req, domains) //nolint:bodyclose
	if err != nil {
		return ListDomainsResponse{}, err
	}

	return *domains, nil
}

// ListWithContext returns the list of all domains
// https://resend.com/docs/api-reference/domains/list-domains
func (s *DomainsSvcImpl) ListWithContext(ctx context.Context) (ListDomainsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List returns the list of all domains
// https://resend.com/docs/api-reference/domains/list-domains
func (s *DomainsSvcImpl) List() (ListDomainsResponse, error) {
	return s.ListWithContext(context.Background())
}

// RemoveWithContext removes a given domain entry by id
// https://resend.com/docs/api-reference/domains/delete-domain
func (s *DomainsSvcImpl) RemoveWithContext(ctx context.Context, domainId string) (bool, error) { //nolint:revive
	path := "domains/" + domainId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return false, ErrFailedToCreateDomainsRemoveRequest
	}

	// Send Request
	_, err = s.client.Perform(req, nil) //nolint:bodyclose
	if err != nil {
		return false, err
	}

	return true, nil
}

// Remove removes a given domain entry by id
// https://resend.com/docs/api-reference/domains/delete-domain
func (s *DomainsSvcImpl) Remove(domainId string) (bool, error) { //nolint:revive
	return s.RemoveWithContext(context.Background(), domainId)
}

// GetWithContext retrieves a domain object
// https://resend.com/docs/api-reference/domains/get-domain
func (s *DomainsSvcImpl) GetWithContext(ctx context.Context, domainId string) (Domain, error) { //nolint:revive
	path := "domains/" + domainId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Domain{}, ErrFailedToCreateDomainsGetRequest
	}

	domain := new(Domain)

	// Send Request
	_, err = s.client.Perform(req, domain) //nolint:bodyclose
	if err != nil {
		return Domain{}, err
	}

	return *domain, nil
}

// Get retrieves a domain object
// https://resend.com/docs/api-reference/domains/get-domain
func (s *DomainsSvcImpl) Get(domainId string) (Domain, error) { //nolint:revive
	return s.GetWithContext(context.Background(), domainId)
}
