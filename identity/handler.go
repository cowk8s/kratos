// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"net/http"

	"github.com/cowk8s/kratos/x"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/herodot"

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

// swagger:route GET /admin/identities identity listIdentities
//
// # List Identities
//
// Lists all [identities](https://www.ory.sh/docs/kratos/concepts/identity-user-model) in the system.
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Security:
//	  oryAccessToken:
//
//	Responses:
//	  200: listIdentities
//	  default: errorGeneric
func (h *Handler) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	includeCredentials := r.URL.Query()["include_credential"]
	var err error
	var declassify []CredentialsType
	for _, v := range includeCredentials {
		tc, ok := ParseCredentialsType(v)
		if ok {
			declassify = append(declassify, tc)
		} else {
			h.r.Writer().WriteError(w, r, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Invalid value `%s` for parameter `include_credential`.", declassify)))
			return
		}
	}

	var orgId uuid.UUID
	if orgIdStr := r.URL.Query().Get("organization_id"); orgIdStr != "" {
		orgId, err = uuid.FromString(r.URL.Query().Get("organization_id"))
		if err != nil {
			h.r.Writer().WriteError(w, r, errors.WithStack(herodot.ErrBadRequest.WithReasonf("Invalid UUID value `%s` for parameter `organization_id`.", r.URL.Query().Get("organization_id"))))
			return
		}
	}

	h.r.Writer().Write(w, r, isam)
}

// Get Identity Parameters
//
// swagger:parameters getIdentity
//
//nolint:deadcode,unused
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type getIdentity struct {
	// ID must be set to the ID of identity you want to get
	//
	// required: true
	// in: path
	ID string `json:"id"`

	// Include Credentials in Response
	//
	// Include any credential, for example `password` or `oidc`, in the response. When set to `oidc`, This will return
	// the initial OAuth 2.0 Access Token, OAuth 2.0 Refresh Token and the OpenID Connect ID Token if available.
	//
	// required: false
	// in: query
	DeclassifyCredentials []CredentialsType `json:"include_credential"`
}
