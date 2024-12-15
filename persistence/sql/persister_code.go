// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"

	"github.com/gofrs/uuid"
)

type oneTimeCodeProvider interface {
	GetID() uuid.UUID
	Validate() error
	TableName(ctx context.Context) string
	GetHMACCode() string
}

type codeOptions struct {
	IdentityID *uuid.UUID
}

type codeOption func(o *codeOptions)

func withCheckIdentityID(id uuid.UUID) codeOption {
	return func(o *codeOptions) {
		o.IdentityID = &id
	}
}
