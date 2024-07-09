// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cowk8s/kratos/text"

	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/x/stringsx"

	"github.com/samber/lo"
)

type secureRedirectOptions struct {
	allowlist       []url.URL
	defaultReturnTo *url.URL
	returnTo        string
	sourceURL       string
}

type SecureRedirectOption func(*secureRedirectOptions)

// SecureRedirectToIsAllowedHost validates if the redirect_to param is allowed for a given wildcard
func SecureRedirectToIsAllowedHost(returnTo *url.URL, allowed url.URL) bool {
	if allowed.Host != "" && allowed.Host[:1] == "*" {
		return strings.HasSuffix(strings.ToLower(returnTo.Host), strings.ToLower(allowed.Host)[1:])
	}
	return strings.EqualFold(allowed.Host, returnTo.Host)
}

// SecureRedirectUseSourceURL uses the given source URL (checks the `?return_to` value)
// instead of r.URL.
func SecureRedirectUseSourceURL(source string) SecureRedirectOption {
	return func(o *secureRedirectOptions) {
		o.sourceURL = source
	}
}

// SecureRedirectTo implements a HTTP redirector who mitigates open redirect vulnerabilities by
// working with allow lists.
func SecureRedirectTo(r *http.Request, defaultReturnTo *url.URL, opts ...SecureRedirectOption) (returnTo *url.URL, err error) {
	o := &secureRedirectOptions{defaultReturnTo: defaultReturnTo}
	for _, opt := range opts {
		opt(o)
	}

	if len(o.allowlist) == 0 {
		return o.defaultReturnTo, nil
	}

	source := RequestURL(r)
	if o.sourceURL != "" {
		source, err = url.ParseRequestURI(o.sourceURL)
		if err != nil {
			return nil, errors.WithStack(herodot.ErrInternalServerError.WithWrap(err).WithReasonf("Unable to parse the original request URL: %s", err))
		}
	}

	rawReturnTo := stringsx.Coalesce(o.returnTo, source.Query().Get("return_to"))
	if rawReturnTo == "" {
		return o.defaultReturnTo, nil
	}

	returnTo, err = url.Parse(rawReturnTo)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrBadRequest.WithWrap(err).WithReasonf("Unable to parse the return_to query parameter as an URL: %s", err))
	}

	returnTo.Host = stringsx.Coalesce(returnTo.Host, o.defaultReturnTo.Host)
	returnTo.Scheme = stringsx.Coalesce(returnTo.Scheme, o.defaultReturnTo.Scheme)

	for _, allowed := range o.allowlist {
		if strings.EqualFold(allowed.Scheme, returnTo.Scheme) &&
			SecureRedirectToIsAllowedHost(returnTo, allowed) &&
			strings.HasPrefix(
				stringsx.Coalesce(returnTo.Path, "/"),
				stringsx.Coalesce(allowed.Path, "/")) {
			return returnTo, nil
		}
	}

	return nil, errors.WithStack(herodot.ErrBadRequest.
		WithID(text.ErrIDRedirectURLNotAllowed).
		WithReasonf("Requested return_to URL %q is not allowed.", returnTo).
		WithDebugf("Allowed domains are: %v", strings.Join(lo.Map(o.allowlist, func(u url.URL, _ int) string {
			return u.String()
		}), ", ")))
}
