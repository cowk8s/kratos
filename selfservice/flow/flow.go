// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/cowk8s/kratos/driver/config"
	"github.com/cowk8s/kratos/x"
)

type Flow interface {
	GetID() uuid.UUID
	GetType() Type
	GetRequestURL() string
}

type FlowWithRedirect interface {
	SecureRedirectToOpts(ctx context.Context, cfg config.Provider) (opts []x.SecureRedirectOption)
}
