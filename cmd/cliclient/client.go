// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cliclient

import (
	"net/http"

	"github.com/spf13/cobra"
)

type ContextKey int

type ClientContext struct {
	Endpoint   string
	HTTPClient *http.Client
}

func NewClient(cmd *cobra.Command)
