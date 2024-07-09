// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hydra

import (
	"net/http"

	"github.com/cowk8s/kratos/driver/config"
	"github.com/cowk8s/kratos/x"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/x/sqlxx"
)

type (
	hydraDependencies interface {
		config.Provider
		x.HTTPClientProvider
	}

	DefaultHydra struct {
		d hydraDependencies
	}
)

func NewDefaultHydra(d hydraDependencies) *DefaultHydra {
	return &DefaultHydra{
		d: d,
	}
}

func GetLoginChallengeID(conf *config.Config, r *http.Request) (sqlxx.NullString, error) {
	if !r.URL.Query().Has("login_challenge") {
		return "", nil
	} else if conf.OAuth2ProviderURL(r.Context()) == nil {
		return "", errors.WithStack(herodot.ErrInternalServerError.WithReason("refusing to parse login_challenge query parameter because " + config.ViperKeyOAuth2ProviderURL + " is invalid or unset"))
	}

	loginChallenge := r.URL.Query().Get("login_challenge")
	if loginChallenge == "" {
		return "", errors.WithStack(herodot.ErrBadRequest.WithReason("the login_challenge parameter is present but empty"))
	}

	return sqlxx.NullString(loginChallenge), nil
}
