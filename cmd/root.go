// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cowk8s/kratos/cmd/identities"
)

func NewRootCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "kratos",
	}

	cmd.AddCommand(identities.NewGetCmd())
	return cmd
}

func Execute() int {
	c := NewRootCmd()

	if err := c.Execute(); err != nil {
		return 1
	}
	return 0
}
