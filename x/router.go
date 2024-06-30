// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

type RouterAdmin struct {
	*httprouter.Router
}

func NewRouterAdmin() *RouterAdmin {
	return &RouterAdmin{
		Router: httprouter.New(),
	}
}

func (r *RouterAdmin) GET(publicPath string, handle httprouter.Handle) {
	r.Router.GET(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) HEAD(publicPath string, handle httprouter.Handle) {
	r.Router.HEAD(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) POST(publicPath string, handle httprouter.Handle) {
	r.Router.POST(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) PUT(publicPath string, handle httprouter.Handle) {
	r.Router.PUT(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) PATCH(publicPath string, handle httprouter.Handle) {
	r.Router.PATCH(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) DELETE(publicPath string, handle httprouter.Handle) {
	r.Router.DELETE(path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) Handle(method, publicPath string, handle httprouter.Handle) {
	r.Router.Handle(method, path.Join(AdminPrefix, publicPath), NoCacheHandle(handle))
}

func (r *RouterAdmin) HandlerFunc(method, publicPath string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, path.Join(AdminPrefix, publicPath), NoCacheHandlerFunc(handler))
}

func (r *RouterAdmin) Handler(method, publicPath string, handler http.Handler) {
	r.Router.Handler(method, path.Join(AdminPrefix, publicPath), NoCacheHandler(handler))
}

func (r *RouterAdmin) Lookup(method, publicPath string) {
	r.Router.Lookup(method, path.Join(AdminPrefix, publicPath))
}
