// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import "github.com/ory/x/sqlxx"

// CredentialsConfig is the struct that is being used as part of the identity credentials.
type CredentialsLookupConfig struct {
	// List of recovery codes
	RecoveryCodes []RecoveryCode `json:"recovery_codes"`
}

type RecoveryCode struct {
	// A recovery code
	Code string `json:"code"`

	// UsedAt indicates whether and when a recovery code was used.
	UsedAt sqlxx.NullTime `json:"used_at,omitempty"`
}
