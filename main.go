// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// package main is the entry point for kratos.
package main

import (
	"os"

	"github.com/cowk8s/kratos/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
