// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/x/httpx"
	"github.com/ory/x/stringsx"
)

func RequestURL(r *http.Request) *url.URL {
	source := *r.URL
	source.Host = stringsx.Coalesce(source.Host, r.Header.Get("X-Forwarded-Host"), r.Host)

	if proto := r.Header.Get("X-Forwarded-Proto"); len(proto) > 0 {
		source.Scheme = proto
	}

	if source.Scheme == "" {
		source.Scheme = "https"
		if r.TLS == nil {
			source.Scheme = "http"
		}
	}

	return &source
}

type HTTPClientProvider interface {
	HTTPClient(context.Context, ...httpx.ResilientOptions) *retryablehttp.Client
}
