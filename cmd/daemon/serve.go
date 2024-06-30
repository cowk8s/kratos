// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package daemon

import (
	stdctx "context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/ory/x/otelx/semconv"

	"github.com/pkg/errors"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"

	"github.com/cowk8s/kratos/driver"
	"github.com/cowk8s/kratos/x"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/ory/graceful"
	"github.com/ory/x/healthx"
	"github.com/ory/x/networkx"
	prometheus "github.com/ory/x/prometheusx"
	"github.com/ory/x/reqlog"

	"github.com/ory/x/servicelocatorx"
)

type options struct {
	ctx stdctx.Context
}

func NewOptions(ctx stdctx.Context, opts []Option) *options {
	o := new(options)
	o.ctx = ctx
	for _, f := range opts {
		f(o)
	}
	return o
}

type Option func(*options)

func WithContext(ctx stdctx.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func init() {
	graceful.DefaultShutdownTimeout = 120 * time.Second
}

func serveAdmin(r driver.Registry, cmd *cobra.Command, eg *errgroup.Group, slOpts *servicelocatorx.Options, opts []Option) {
	modifiers := NewOptions(cmd.Context(), opts)
	ctx := modifiers.ctx

	c := r.Config()
	l := r.Logger()
	n := negroni.New()

	for _, mw := range slOpts.HTTPMiddlewares() {
		n.UseFunc(mw)
	}

	adminLogger := reqlog.NewMiddlewareFromLogger(
		l,
		"admin#"+c.SelfPublicURL(ctx).String(),
	)

	if r.Config().DisableAdminHealthRequestLog(ctx) {
		adminLogger.ExcludePaths(x.AdminPrefix+healthx.AliveCheckPath, x.AdminPrefix+healthx.ReadyCheckPath, x.AdminPrefix+prometheus.MetricsPrometheusPath)
	}
	n.UseFunc(semconv.Middleware)
	n.Use(adminLogger)

	router := x.NewRouterAdmin()
	r.RegisterAdminRoutes(ctx, router)

	n.UseHandler(http.MaxBytesHandler(router, 5*1024*1024 /* 5 MB */))
	certs := c.GetTLSCertificatesForAdmin(ctx)

	var handler http.Handler = n

	//#nosec G112 -- the correct settings are set by graceful.WithDefaults
	server := graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         &tls.Config{GetCertificate: certs, MinVersion: tls.VersionTLS12},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       600 * time.Second,
	})

	addr := c.AdminListenOn(ctx)

	eg.Go(func() error {
		l.Printf("Starting the admin httpd on: %s", addr)
		if err := graceful.GracefulContext(ctx, func() error {
			listener, err := networkx.MakeListener(addr, c.AdminSocketPermission(ctx))
			if err != nil {
				return err
			}

			if certs == nil {
				return server.Serve(listener)
			}
			return server.ServeTLS(listener, "", "")
		}, server.Shutdown); err != nil {
			if !errors.Is(err, context.Canceled) {
				l.Errorf("Failed to gracefully shutdown admin httpd: %s", err)
				return err
			}
		}
		l.Println("Admin httpd was shutdown gracefully")
		return nil
	})
}

func ServeAll(d driver.Registry, cmd *cobra.Command, slOpts *servicelocatorx.Options, opts []Option) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		mods := NewOptions(cmd.Context(), opts)
		ctx := mods.ctx

		g, ctx := errgroup.WithContext(ctx)
		cmd.SetContext(ctx)
		opts = append(opts, WithContext(ctx))

		serveAdmin(d, cmd, g, slOpts, opts)

		return g.Wait()
	}
}
