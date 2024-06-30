// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identities

import (
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
	}
	cmd.AddCommand(NewGetIdentityCmd())
	return cmd
}

func NewGetIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity [id-1] [id-2] [id-n]",
		Short: "Get one or more identities by their ID(s)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
