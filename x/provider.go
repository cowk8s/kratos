// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
)

type LoggingProvider interface {
	Logger() *logrusx.Logger
	Audit() *logrusx.Logger
}

type TracingProvider interface {
	Tracer(ctx context.Context) *otelx.Tracer
}
