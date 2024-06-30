// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"io"
)

func New(ctx context.Context, stdOutOrErr io.Writer) error {
}

func NewWithoutInit(ctx context.Context, stdOutOrErr io.Writer) error {
}
