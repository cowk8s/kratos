// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import "github.com/ory/x/sqlxx"

type (
	Expandable  = sqlxx.Expandable
	Expandables = sqlxx.Expandables
)

const (
	ExpandFieldVerifiableAddresses Expandable = "VerifiableAddresses"
	ExpandFieldRecoveryAddresses   Expandable = "RecoveryAddresses"
	ExpandFieldCredentials         Expandable = "Credentials"
)

// ExpandNothing expands nothing
var ExpandNothing = Expandables{}
