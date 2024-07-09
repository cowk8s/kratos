// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package registration

import (
	"net/http"
	"time"

	"github.com/cowk8s/kratos/driver/config"
	"github.com/cowk8s/kratos/hydra"
	"github.com/cowk8s/kratos/selfservice/flow"
	"github.com/cowk8s/kratos/x"
	"github.com/gofrs/uuid"

	"github.com/ory/x/sqlxx"
)

type Flow struct {
	ID uuid.UUID `json:"id" faker:"-" db:"id"`

	OAuth2LoginChallenge sqlxx.NullString `json:"oauth2_login_challenge,omitempty" faker:"-" db:"oauth2_login_challenge_data"`

	Type flow.Type `json:"type" db:"type" faker:"flow_type"`

	ExpiresAt time.Time `json:"expires_at" faker:"time_type" db:"expires_at"`

	IssuedAt time.Time `json:"issued_at" faker:"time_type" db:"issued_at"`

	RequestURL string `json:"request_url" faker:"url" db:"request_url"`

	UI *container
}

var _ flow.Flow = new(Flow)

func NewFlow(conf *config.Config, exp time.Duration, csrf string, r *http.Request, flowType flow.Type) (*Flow, error) {
	now := time.Now().UTC()
	id := x.NewUUID()
	requestURL := x.RequestURL(r).String()
	_, err := x.SecureRedirectTo(r,
		conf.SelfServiceBrowserDefaultReturnTo(r.Context()),
		x.SecureRedirectUseSourceURL(requestURL),
	)
	if err != nil {
		return nil, err
	}

	hlc, err := hydra.GetLoginChallengeID(conf, r)
	if err != nil {
		return nil, err
	}

	return &Flow{
		ID:                   id,
		OAuth2LoginChallenge: hlc,
		ExpiresAt:            now.Add(exp),
		IssuedAt:             now,
		RequestURL:           requestURL,
	}, nil
}

func (f Flow) GetID() uuid.UUID {
	return f.ID
}

func (f Flow) GetType() flow.Type {
	return f.Type
}
