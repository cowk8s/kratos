// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"context"

	"github.com/cowk8s/kratos/driver/config"

	"github.com/ory/herodot"
)

var ErrNoAALAvailable = herodot.ErrForbidden.WithReasonf("Unable to detect available authentication methods. Perform account recovery or contact support.")

type (
	managerHTTPDependencies interface {
		config.Provider
	}
	ManagerHttp struct {
		cookieName func(ctx context.Context) string
		r          managerHTTPDependencies
	}
)

type options struct {
	requestURL string
	upsertAAL  bool
}

type ManagerOptions func(*options)
