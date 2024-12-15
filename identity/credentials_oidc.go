// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

// CredentialsOIDC is contains the configuration for credentials of the type oidc.
//
// swagger:model identityCredentialsOidc
type CredentialsOIDC struct {
	Providers []CredentialsOIDCProvider `json:"providers"`
}

// CredentialsOIDCProvider is contains a specific OpenID COnnect credential for a particular connection (e.g. Google).
//
// swagger:model identityCredentialsOidcProvider
type CredentialsOIDCProvider struct {
	Subject             string `json:"subject"`
	Provider            string `json:"provider"`
	InitialIDToken      string `json:"initial_id_token"`
	InitialAccessToken  string `json:"initial_access_token"`
	InitialRefreshToken string `json:"initial_refresh_token"`
	Organization        string `json:"organization,omitempty"`
}

// swagger:ignore
type CredentialsOIDCEncryptedTokens struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
}
