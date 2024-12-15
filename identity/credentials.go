// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import "github.com/gofrs/uuid"

// Authenticator Assurance Level (AAL)
//
// The authenticator assurance level can be one of "aal1", "aal2", or "aal3". A higher number means that it is harder
// for an attacker to compromise the account.
//
// Generally, "aal1" implies that one authentication factor was used while AAL2 implies that two factors (e.g.
// password + TOTP) have been used.
//
// To learn more about these levels please head over to: https://www.ory.sh/kratos/docs/concepts/credentials
//
// swagger:model authenticatorAssuranceLevel
type AuthenticatorAssuranceLevel string

const (
	NoAuthenticatorAssuranceLevel AuthenticatorAssuranceLevel = "aal0"
	AuthenticatorAssuranceLevel1  AuthenticatorAssuranceLevel = "aal1"
	AuthenticatorAssuranceLevel2  AuthenticatorAssuranceLevel = "aal2"
)

// CredentialsType  represents several different credential types, like password credentials, passwordless credentials,
// and so on.
//
// swagger:enum CredentialsType
type CredentialsType string

// Please make sure to add all of these values to the test that ensures they are created during migration
const (
	CredentialsTypePassword CredentialsType = "password"
	CredentialsTypeOIDC     CredentialsType = "oidc"
	CredentialsTypeTOTP     CredentialsType = "totp"
	CredentialsTypeLookup   CredentialsType = "lookup_secret"
	CredentialsTypeWebAuthn CredentialsType = "webauthn"
	CredentialsTypeCodeAuth CredentialsType = "code"
	CredentialsTypePasskey  CredentialsType = "passkey"
	CredentialsTypeProfile  CredentialsType = "profile"
)

func (c CredentialsType) String() string {
	return string(c)
}

const (
	// CredentialsTypeRecoveryLink is a special credential type linked to the link strategy (recovery flow).
	// It is not used within the credentials object itself.
	CredentialsTypeRecoveryLink CredentialsType = "link_recovery"
	CredentialsTypeRecoveryCode CredentialsType = "code_recovery"
)

func ParseCredentialsType(in string) (CredentialsType, bool) {
	for _, t := range []CredentialsType{
		CredentialsTypePassword,
	} {
		if t.String() == in {
			return t, true
		}
	}
	return "", false
}

// Credentials represents a specific credential type
//
// swagger:model identityCredentials
type Credentials struct {
	ID uuid.UUID `json:"-" db:"id"`

	// Type discriminates between different types of credentials.
	Type                     CredentialsType `json:"type" db:"-"`
	IdentityCredentialTypeID uuid.UUID       `json:"-" db:"identity_credential_type_id"`
}
