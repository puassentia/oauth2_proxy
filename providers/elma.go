package providers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type ELMAProvider struct {
	*ProviderData
}

func NewELMAProvider(p *ProviderData) *ELMAProvider {
	p.ProviderName = "ELMA"
	if p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{Scheme: "https",
			Host: "elma.elewise.com",
			Path: "/oauth/oauth/authorize"}
	}
	if p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{Scheme: "https",
			Host: "elma.elewise.com",
			Path: "/oauth/oauth/token"}
	}
	if p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{Scheme: "https",
			Host: "elma.elewise.com",
			Path: "/oauth/oauth/userinfo"}
	}
	if p.ValidateURL.String() == "" {
		p.ValidateURL = p.ProfileURL
	}
	if p.Scope == "" {
		p.Scope = "r_emailaddress r_basicprofile"
	}
	return &ELMAProvider{ProviderData: p}
}

func getELMAHeader(access_token string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("x-li-format", "json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return header
}

func (p *ELMAProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String()+"?format=json", nil)
	if err != nil {
		return "", err
	}
	req.Header = getELMAHeader(s.AccessToken)

	json, err := api.Request(req)
	if err != nil {
		return "", err
	}

	email, err := json.String()
	if err != nil {
		return "", err
	}
	return email, nil
}

func (p *ELMAProvider) ValidateSessionState(s *SessionState) bool {
	return validateToken(p, s.AccessToken, getELMAHeader(s.AccessToken))
}
