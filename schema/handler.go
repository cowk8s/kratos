// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/cowk8s/kratos/x"
)

type (
	handlerDependencies interface {
		x.CSRFProvider
	}
	Handler struct {
		r handlerDependencies
	}
	HandlerProvider interface {
		SchemaHandler() *Handler
	}
)

func NewHandler(r handlerDependencies) *Handler {
	return &Handler{r: r}
}

const SchemasPath string = "schemas"

func (h *Handler) RegisterPublicRoutes(public *x.RouterPublic) {
	h.r.C
}
