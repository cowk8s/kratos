// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"net/http"

	"github.com/cowk8s/kratos/x"
	"github.com/julienschmidt/httprouter"

	"github.com/ory/x/decoderx"
)

const (
	RouteCollection = "/identities"
)

type (
	handlerDependencies interface{}
	Handler             struct {
		r  handlerDependencies
		dx *decoderx.HTTP
	}
)

func NewHandler(r handlerDependencies) *Handler {
	return &Handler{
		r:  r,
		dx: decoderx.NewHTTP(),
	}
}

func (h *Handler) RegisterAdminRoutes(admin *x.RouterAdmin) {
	admin.GET(RouteCollection, h.list)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	include
}
