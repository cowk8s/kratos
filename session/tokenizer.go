// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"time"

	"github.com/dgraph-io/ristretto"

	"github.com/cowk8s/kratos/driver/config"
	"github.com/cowk8s/kratos/x"

	"github.com/ory/x/jsonnetsecure"
)

type (
	tokenizerDependencies interface {
		jsonnetsecure.VMProvider
		x.TracingProvider
		x.HTTPClientProvider
		config.Provider
		x.JWKSFetchProvider
	}
	Tokenizer struct {
		r       tokenizerDependencies
		nowFunc func() time.Time
		cache   *ristretto.Cache[[]byte, []byte]
	}
	TokenizerProvider interface {
		SessionTokenizer() *Tokenizer
	}
)
