// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package serve

import (
	"github.com/spf13/cobra"

	"github.com/cowk8s/kratos/cmd/daemon"
)

func NewServeCmd() (serveCmd *cobra.Command) {
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the Ory Kratos server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			return daemon.ServeAll(cmd, nil)(cmd, args)
		},
	}
}
