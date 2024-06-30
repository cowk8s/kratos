// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/tlsx"
)

const (
	ViperKeyAdminSocketOwner             = "serve.admin.socket.owner"
	ViperKeyAdminSocketGroup             = "serve.admin.socket.group"
	ViperKeyAdminSocketMode              = "serve.admin.socket.mode"
	ViperKeyPublicTLSCertBase64          = "serve.public.tls.cert.base64"
	ViperKeyPublicTLSKeyBase64           = "serve.public.tls.key.base64"
	ViperKeyPublicTLSCertPath            = "serve.public.tls.cert.path"
	ViperKeyPublicTLSKeyPath             = "serve.public.tls.key.path"
	ViperKeyDisableAdminHealthRequestLog = "serve.admin.request_log.disable_for_health"

	ViperKeyAdminTLSCertBase64 = "serve.admin.tls.cert.base64"
	ViperKeyAdminTLSKeyBase64  = "serve.admin.tls.key.base64"
	ViperKeyAdminTLSCertPath   = "serve.admin.tls.cert.path"
	ViperKeyAdminTLSKeyPath    = "serve.admin.tls.key.path"

	ViperKeyPublicBaseURL = "serve.public.base_url"
	ViperKeyPublicPort    = "serve.public.port"
	ViperKeyPublicHost    = "serve.public.host"
)

type (
	Config struct {
		l           *logrusx.Logger
		p           *configx.Provider
		c           contextx.Contextualizer
		stdOutOrErr io.Writer
	}

	Provider interface {
		Config() *Config
	}
)

func (p *Config) DSN(ctx context.Context) string {
	return ""
}

func (p *Config) listenOn(ctx context.Context, key string) string {
	fb := 4433
	if key == "admin" {
		fb = 4434
	}

	pp := p.GetProvider(ctx)
	port := pp.IntF("serve."+key+".port", fb)
	if port < 1 {
		p.l.Fatalf("serve.%s.port can not be zero or negative", key)
	}

	return configx.GetAddress(pp.String("serve."+key+".host"), port)
}

func (p *Config) AdminListenOn(ctx context.Context) string {
	return p.listenOn(ctx, "admin")
}

func (p *Config) PublicListenOn(ctx context.Context) string {
	return p.listenOn(ctx, "public")
}

func (p *Config) AdminSocketPermission(ctx context.Context) *configx.UnixPermission {
	pp := p.GetProvider(ctx)
	return &configx.UnixPermission{
		Owner: pp.String(ViperKeyAdminSocketOwner),
		Group: pp.String(ViperKeyAdminSocketGroup),
		Mode:  os.FileMode(pp.IntF(ViperKeyAdminSocketMode, 0o755)),
	}
}

func (p *Config) guessBaseURL(ctx context.Context, keyHost, keyPort string, defaultPort int) *url.URL {
	port := p.GetProvider(ctx).IntF(keyPort, defaultPort)

	host := p.GetProvider(ctx).String(keyHost)
	if host == "0.0.0.0" || len(host) == 0 {
		var err error
		host, err = os.Hostname()
		if err != nil {
			p.l.WithError(err).Warn("Unable to get hostname from system, falling back to 127.0.0.1.")
			host = "127.0.0.1"
		}
	}

	guess := url.URL{Host: fmt.Sprintf("%s:%d", host, port), Scheme: "https", Path: "/"}
	if p.IsInsecureDevMode(ctx) {
		guess.Scheme = "http"
	}

	return &guess
}

func (p *Config) baseURL(ctx context.Context, keyURL, keyHost, keyPort string, defaultPort int) *url.URL {
	switch t := p.GetProvider(ctx).Get(keyURL).(type) {
	case *url.URL:
		return t
	case url.URL:
		return &t
	case string:
		parsed, err := url.ParseRequestURI(t)
		if err != nil {
			p.l.WithError(err).Errorf("Configuration key %s is not a valid URL. Falling back to optimistically guessing the server's base URL. Please set a value to avoid problems with redirects and cookies.", keyURL)
			return p.guessBaseURL(ctx, keyHost, keyPort, defaultPort)
		}
		return parsed
	}

	p.l.Warnf("Configuration key %s was left empty. Optimistically guessing the server's base URL. Please set a value to avoid problems with redirects and cookies.", keyURL)
	return p.guessBaseURL(ctx, keyHost, keyPort, defaultPort)
}

func (p *Config) SelfPublicURL(ctx context.Context) *url.URL {
	return p.baseURL(ctx, ViperKeyPublicBaseURL, ViperKeyPublicHost, ViperKeyPublicPort, 4433)
}

func (p *Config) DisableAdminHealthRequestLog(ctx context.Context) bool {
	return p.GetProvider(ctx).Bool(ViperKeyDisableAdminHealthRequestLog)
}

func (p *Config) IsInsecureDevMode(ctx context.Context) bool {
	return p.GetProvider(ctx).Bool("dev")
}

type CertFunc = func(*tls.ClientHelloInfo) (*tls.Certificate, error)

func (p *Config) GetTLSCertificatesForPublic(ctx context.Context) CertFunc {
	return p.getTLSCertificates(
		ctx,
		"public",
		p.GetProvider(ctx).String(ViperKeyPublicTLSCertBase64),
		p.GetProvider(ctx).String(ViperKeyPublicTLSKeyBase64),
		p.GetProvider(ctx).String(ViperKeyPublicTLSCertPath),
		p.GetProvider(ctx).String(ViperKeyPublicTLSKeyPath),
	)
}

func (p *Config) GetTLSCertificatesForAdmin(ctx context.Context) CertFunc {
	return p.getTLSCertificates(
		ctx,
		"admin",
		p.GetProvider(ctx).String(ViperKeyAdminTLSCertBase64),
		p.GetProvider(ctx).String(ViperKeyAdminTLSKeyBase64),
		p.GetProvider(ctx).String(ViperKeyAdminTLSCertPath),
		p.GetProvider(ctx).String(ViperKeyAdminTLSKeyPath),
	)
}

func (p *Config) getTLSCertificates(ctx context.Context, daemon, certBase64, keyBase64, certPath, keyPath string) CertFunc {
	if certBase64 != "" && keyBase64 != "" {
		cert, err := tlsx.CertificateFromBase64(certBase64, keyBase64)
		if err != nil {
			p.l.WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
			return nil // reachable in unit tests when Fatalf is hooked
		}
		p.l.Infof("Setting up HTTPS for %s", daemon)
		return func(*tls.ClientHelloInfo) (*tls.Certificate, error) { return &cert, nil }
	}
	if certPath != "" && keyPath != "" {
		errs := make(chan error, 1)
		getCert, err := tlsx.GetCertificate(ctx, certPath, keyPath, errs)
		if err != nil {
			p.l.WithError(err).Fatalf("Unable to load HTTPS TLS Certificate")
			return nil // reachable in unit tests when Fatalf is hooked
		}
		go func() {
			for err := range errs {
				p.l.WithError(err).Error("Failed to reload TLS certificates, using previous certificates")
			}
		}()
		p.l.Infof("Setting up HTTPS for %s (automatic certificate reloading active)", daemon)
		return getCert
	}
	p.l.Infof("TLS has not been configured for %s, skipping", daemon)
	return nil
}

func (p *Config) GetProvider(ctx context.Context) *configx.Provider {
	return p.c.Config(ctx, p.p)
}
