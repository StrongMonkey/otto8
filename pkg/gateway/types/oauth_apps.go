package types

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/otto8-ai/otto8/apiclient/types"
)

const (
	SlackAuthorizeURL = "https://slack.com/oauth/v2/authorize"
	SlackTokenURL     = "https://slack.com/api/oauth.v2.access"

	NotionAuthorizeURL = "https://api.notion.com/v1/oauth/authorize"
	NotionTokenURL     = "https://api.notion.com/v1/oauth/token"

	HubSpotAuthorizeURL = "https://app.hubspot.com/oauth/authorize"
	HubSpotTokenURL     = "https://api.hubapi.com/oauth/v1/token"

	GoogleAuthorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL     = "https://oauth2.googleapis.com/token"
)

type OAuthAppTypeConfig struct {
	DisplayName string            `json:"displayName"`
	Parameters  map[string]string `json:"parameters"`
}

func SupportedOAuthAppTypeConfigs() map[types.OAuthAppType]OAuthAppTypeConfig {
	return map[types.OAuthAppType]OAuthAppTypeConfig{
		types.OAuthAppTypeMicrosoft365: {
			DisplayName: "Microsoft 365",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
				"tenantID":     "Tenant ID",
			},
		},
		types.OAuthAppTypeSlack: {
			DisplayName: "Slack",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
			},
		},
		types.OAuthAppTypeNotion: {
			DisplayName: "Notion",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
			},
		},
		types.OAuthAppTypeHubSpot: {
			DisplayName: "HubSpot",
			Parameters: map[string]string{
				"name":          "App Name",
				"clientID":      "Client ID",
				"clientSecret":  "Client Secret",
				"appID":         "App ID",
				"optionalScope": "Optional Scope",
			},
		},
		types.OAuthAppTypeGitHub: {
			DisplayName: "GitHub",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
				"authURL":      "Authorization URL",
				"tokenURL":     "Token URL",
			},
		},
		types.OAuthAppTypeGoogle: {
			DisplayName: "Google",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
			},
		},
		types.OAuthAppTypeCustom: {
			DisplayName: "Custom",
			Parameters: map[string]string{
				"name":         "App Name",
				"clientID":     "Client ID",
				"clientSecret": "Client Secret",
				"authURL":      "Authorization URL",
				"tokenURL":     "Token URL",
			},
		},
	}
}

func ValidateAndSetDefaultsOAuthAppManifest(r *types.OAuthAppManifest) error {
	var errs []error
	// Check for missing fields.
	if r.ClientID == "" {
		errs = append(errs, fmt.Errorf("missing clientID"))
	} else if r.Type == types.OAuthAppTypeCustom && r.AuthURL == "" {
		errs = append(errs, fmt.Errorf("missing authURL"))
	} else if r.Type == types.OAuthAppTypeCustom && r.TokenURL == "" {
		errs = append(errs, fmt.Errorf("missing tokenURL"))
	} else if r.Type == types.OAuthAppTypeMicrosoft365 && r.TenantID == "" {
		errs = append(errs, fmt.Errorf("missing tenantID"))
	} else if r.Type == types.OAuthAppTypeHubSpot && r.AppID == "" {
		errs = append(errs, fmt.Errorf("missing appID"))
	}

	switch r.Type {
	case types.OAuthAppTypeMicrosoft365:
		r.AuthURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", r.TenantID)
		r.TokenURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", r.TenantID)
	case types.OAuthAppTypeSlack:
		r.AuthURL = SlackAuthorizeURL
		r.TokenURL = SlackTokenURL
	case types.OAuthAppTypeNotion:
		r.AuthURL = NotionAuthorizeURL
		r.TokenURL = NotionTokenURL
	case types.OAuthAppTypeHubSpot:
		r.AuthURL = HubSpotAuthorizeURL
		r.TokenURL = HubSpotTokenURL
	case types.OAuthAppTypeGoogle:
		r.AuthURL = GoogleAuthorizeURL
		r.TokenURL = GoogleTokenURL
	}

	// Validate URLs.
	if _, err := url.Parse(r.AuthURL); err != nil {
		errs = append(errs, fmt.Errorf("invalid authURL: %w", err))
	} else if _, err := url.Parse(r.TokenURL); err != nil {
		errs = append(errs, fmt.Errorf("invalid tokenURL: %w", err))
	}

	return errors.Join(errs...)
}

func MergeOAuthAppManifests(r, other *types.OAuthAppManifest) *types.OAuthAppManifest {
	retVal := *r

	if other.RefName != "" {
		retVal.RefName = other.RefName
	}
	if other.ClientID != "" {
		retVal.ClientID = other.ClientID
	}
	if other.AuthURL != "" {
		retVal.AuthURL = other.AuthURL
	}
	if other.TokenURL != "" {
		retVal.TokenURL = other.TokenURL
	}
	if other.Type != "" {
		retVal.Type = other.Type
	}
	if other.TenantID != "" {
		retVal.TenantID = other.TenantID
	}

	return &retVal
}

// OAuthTokenResponse represents a response from the /token endpoint on an OAuth server.
// These do not get stored in the database.
type OAuthTokenResponse struct {
	State        string `json:"state"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Ok           bool   `json:"ok"`
	Error        string `json:"error"`
	CreatedAt    time.Time
}

type GoogleOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type SlackOAuthTokenResponse struct {
	Ok         bool   `json:"ok"`
	Error      string `json:"error"`
	AuthedUser struct {
		ID          string `json:"id"`
		Scope       string `json:"scope"`
		AccessToken string `json:"access_token"`
	} `json:"authed_user"`
}

type OAuthTokenRequestChallenge struct {
	State     string    `json:"state" gorm:"primaryKey"`
	Challenge string    `json:"challenge"`
	CreatedAt time.Time `json:"createdAt"`
}
