// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// An Identity's State
//
// The state can either be `active` or `inactive`.
//
// swagger:enum State
type State string

const (
	StateActive   State = "active"
	StateInactive State = "inactive"
)

func (lt State) IsValid() error {
	switch lt {
	case StateActive, StateInactive:
		return nil
	}
	return errors.New("identity state is not valid")
}

// Identity represents an Ory Kratos identity
//
// An [identity](https://www.ory.sh/docs/kratos/concepts/identity-user-model) represents a (human) user in Ory.
//
// swagger:model identity
type Identity struct {
	l *sync.RWMutex `db:"-" faker:"-"`

	// ID is the identity's unique identifier.
	//
	// The Identity ID can not be changed and can not be chosen. This ensures future
	// compatibility and optimization for distributed stores such as CockroachDB.
	//
	// required: true
	ID uuid.UUID `json:"id" faker:"-" db:"id"`

	// Credentials represents all credentials that can be used for authenticating this identity.
	Credentials map[CredentialsType]Credentials `json:"credentials,omitempty" faker:"-" db:"-"`

	// State is the identity's state.
	//
	// This value has currently no effect.
	State State `json:"state" faker:"-" db:"state"`
}

func (i *Identity) lock() *sync.RWMutex {
	if i.l == nil {
		i.l = new(sync.RWMutex)
	}
	return i.l
}

func (i *Identity) IsActive() bool {
	return i.State == StateActive
}

func (i *Identity) CopyWithoutCredentials() *Identity {
	i.lock().RLock()
	defer i.lock().RUnlock()
	ii := *i
	ii.l = new(sync.RWMutex)
	ii.Credentials = nil
	return &ii
}
