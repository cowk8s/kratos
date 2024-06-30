// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"

	"github.com/ory/x/dbal"

	"github.com/cowk8s/kratos/driver/config"
	"github.com/cowk8s/kratos/x"
)

type Registry interface {
	dbal.Driver

	RegisterAdminRoutes(ctx context.Context, admin *x.RouterAdmin)

	config.Provider

	x.LoggingProvider
}

func NewRegistryFromDSN(ctx context.Context, c *config.Config) error {
	driver, err := dbal.GetDriverFor(c.DSN(ctx))
}
