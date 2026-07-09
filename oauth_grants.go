package resend

import (
	"context"
	"net/http"
)

// OAuthGrantClient represents the OAuth client that a grant was issued to.
type OAuthGrantClient struct {
	Name    string  `json:"name"`
	LogoUri *string `json:"logo_uri"`
}

// OAuthGrant represents an OAuth grant.
type OAuthGrant struct {
	Id            string           `json:"id"`
	ClientId      string           `json:"client_id"`
	Scopes        []string         `json:"scopes"`
	CreatedAt     string           `json:"created_at"`
	RevokedAt     *string          `json:"revoked_at"`
	RevokedReason *string          `json:"revoked_reason"`
	Client        OAuthGrantClient `json:"client"`
}

type ListOAuthGrantsResponse struct {
	Object  string       `json:"object"`
	Data    []OAuthGrant `json:"data"`
	HasMore bool         `json:"has_more"`
}

type RevokeOAuthGrantResponse struct {
	Object        string  `json:"object"`
	Id            string  `json:"id"`
	RevokedAt     *string `json:"revoked_at"`
	RevokedReason *string `json:"revoked_reason"`
}

type OAuthGrantsSvc interface {
	ListWithOptions(ctx context.Context, options *ListOptions) (ListOAuthGrantsResponse, error)
	ListWithContext(ctx context.Context) (ListOAuthGrantsResponse, error)
	List() (ListOAuthGrantsResponse, error)
	RevokeWithContext(ctx context.Context, oauthGrantId string) (RevokeOAuthGrantResponse, error)
	Revoke(oauthGrantId string) (RevokeOAuthGrantResponse, error)
}

type OAuthGrantsSvcImpl struct {
	client *Client
}

// ListWithOptions lists all of the team's OAuth grants with pagination options
// https://resend.com/docs/api-reference/oauth/list-grants
func (s *OAuthGrantsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListOAuthGrantsResponse, error) {
	path := "oauth/grants" + buildPaginationQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListOAuthGrantsResponse{}, ErrFailedToCreateOAuthGrantsListRequest
	}

	grantsResp := new(ListOAuthGrantsResponse)

	_, err = s.client.Perform(req, grantsResp)
	if err != nil {
		return ListOAuthGrantsResponse{}, err
	}

	return *grantsResp, nil
}

// ListWithContext lists all of the team's OAuth grants
// https://resend.com/docs/api-reference/oauth/list-grants
func (s *OAuthGrantsSvcImpl) ListWithContext(ctx context.Context) (ListOAuthGrantsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List lists all of the team's OAuth grants
func (s *OAuthGrantsSvcImpl) List() (ListOAuthGrantsResponse, error) {
	return s.ListWithContext(context.Background())
}

// RevokeWithContext revokes a given OAuth grant by id
// https://resend.com/docs/api-reference/oauth/revoke-grant
func (s *OAuthGrantsSvcImpl) RevokeWithContext(ctx context.Context, oauthGrantId string) (RevokeOAuthGrantResponse, error) {
	path := "oauth/grants/" + oauthGrantId

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RevokeOAuthGrantResponse{}, ErrFailedToCreateOAuthGrantsRevokeRequest
	}

	grantResp := new(RevokeOAuthGrantResponse)

	_, err = s.client.Perform(req, grantResp)
	if err != nil {
		return RevokeOAuthGrantResponse{}, err
	}

	return *grantResp, nil
}

// Revoke revokes a given OAuth grant by id
func (s *OAuthGrantsSvcImpl) Revoke(oauthGrantId string) (RevokeOAuthGrantResponse, error) {
	return s.RevokeWithContext(context.Background(), oauthGrantId)
}
