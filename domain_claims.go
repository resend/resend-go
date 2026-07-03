package resend

import (
	"context"
	"errors"
	"net/http"
)

type DomainClaimStatus = string

const (
	DomainClaimStatusPending    DomainClaimStatus = "pending"
	DomainClaimStatusVerified   DomainClaimStatus = "verified"
	DomainClaimStatusCompleted  DomainClaimStatus = "completed"
	DomainClaimStatusBlocked    DomainClaimStatus = "blocked"
	DomainClaimStatusExpired    DomainClaimStatus = "expired"
	DomainClaimStatusSuperseded DomainClaimStatus = "superseded"
	DomainClaimStatusCanceled   DomainClaimStatus = "canceled"
	DomainClaimStatusFailed     DomainClaimStatus = "failed"
)

type DomainClaimBlockedReason = string

const (
	DomainClaimBlockedReasonGracePeriod            DomainClaimBlockedReason = "grace_period"
	DomainClaimBlockedReasonRecentOwnerActivity    DomainClaimBlockedReason = "recent_owner_activity"
	DomainClaimBlockedReasonPendingScheduledEmails DomainClaimBlockedReason = "pending_scheduled_emails"
)

// DomainClaimRecord is the TXT record to add to your DNS to prove ownership of the claimed domain.
type DomainClaimRecord struct {
	Type  string `json:"type,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Ttl   string `json:"ttl,omitempty"`
}

type DomainClaim struct {
	Object        string             `json:"object,omitempty"`
	Id            string             `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Status        DomainClaimStatus  `json:"status,omitempty"`
	DomainId      string             `json:"domain_id,omitempty"`
	Region        string             `json:"region,omitempty"`
	Record        *DomainClaimRecord `json:"record,omitempty"`
	BlockedReason string             `json:"blocked_reason,omitempty"`
	FailureReason string             `json:"failure_reason,omitempty"`
	CreatedAt     string             `json:"created_at,omitempty"`
	ExpiresAt     string             `json:"expires_at,omitempty"`
}

// CreateDomainClaimRequest contains params for starting a domain claim.
// Uses the same fields as CreateDomainRequest.
type CreateDomainClaimRequest struct {
	Name              string `json:"name"`
	Region            string `json:"region,omitempty"`
	CustomReturnPath  string `json:"custom_return_path,omitempty"`
	TrackingSubdomain string `json:"tracking_subdomain,omitempty"`
	OpenTracking      *bool  `json:"open_tracking,omitempty"`
	ClickTracking     *bool  `json:"click_tracking,omitempty"`
}

type DomainClaimsSvc interface {
	Create(params *CreateDomainClaimRequest) (DomainClaim, error)
	CreateWithContext(ctx context.Context, params *CreateDomainClaimRequest) (DomainClaim, error)
	Get(domainId string) (DomainClaim, error)
	GetWithContext(ctx context.Context, domainId string) (DomainClaim, error)
	Verify(domainId string) (DomainClaim, error)
	VerifyWithContext(ctx context.Context, domainId string) (DomainClaim, error)
}

type DomainClaimsSvcImpl struct {
	client *Client
}

// CreateWithContext starts a claim for a domain that another Resend account has already
// verified, using the given context.
// https://resend.com/docs/api-reference/domains/claim-domain
func (s *DomainClaimsSvcImpl) CreateWithContext(ctx context.Context, params *CreateDomainClaimRequest) (DomainClaim, error) {
	path := "domains/claim"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return DomainClaim{}, errors.New("[ERROR]: Failed to create Domains.Claims.Create request")
	}

	resp := new(DomainClaim)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return DomainClaim{}, err
	}
	return *resp, nil
}

// Create starts a claim for a domain that another Resend account has already verified.
// https://resend.com/docs/api-reference/domains/claim-domain
func (s *DomainClaimsSvcImpl) Create(params *CreateDomainClaimRequest) (DomainClaim, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves the latest claim for the placeholder domain created by the claim,
// using the given context.
// https://resend.com/docs/api-reference/domains/get-domain-claim
func (s *DomainClaimsSvcImpl) GetWithContext(ctx context.Context, domainId string) (DomainClaim, error) {
	path := "domains/" + domainId + "/claim"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return DomainClaim{}, errors.New("[ERROR]: Failed to create Domains.Claims.Get request")
	}

	resp := new(DomainClaim)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return DomainClaim{}, err
	}
	return *resp, nil
}

// Get retrieves the latest claim for the placeholder domain created by the claim.
// https://resend.com/docs/api-reference/domains/get-domain-claim
func (s *DomainClaimsSvcImpl) Get(domainId string) (DomainClaim, error) {
	return s.GetWithContext(context.Background(), domainId)
}

// VerifyWithContext triggers asynchronous DNS verification and ownership transfer for a
// domain claim, using the given context. The claim stays "pending" while verification runs;
// poll GetWithContext for status.
// https://resend.com/docs/api-reference/domains/verify-domain-claim
func (s *DomainClaimsSvcImpl) VerifyWithContext(ctx context.Context, domainId string) (DomainClaim, error) {
	path := "domains/" + domainId + "/claim/verify"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return DomainClaim{}, errors.New("[ERROR]: Failed to create Domains.Claims.Verify request")
	}

	resp := new(DomainClaim)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return DomainClaim{}, err
	}
	return *resp, nil
}

// Verify triggers asynchronous DNS verification and ownership transfer for a domain claim.
// The claim stays "pending" while verification runs; poll Get for status.
// https://resend.com/docs/api-reference/domains/verify-domain-claim
func (s *DomainClaimsSvcImpl) Verify(domainId string) (DomainClaim, error) {
	return s.VerifyWithContext(context.Background(), domainId)
}
