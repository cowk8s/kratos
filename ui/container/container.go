// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/cowk8s/kratos/text"
	"github.com/cowk8s/kratos/ui/node"

	"github.com/ory/x/decoderx"
)

var (
	decoder             = decoderx.NewHTTP()
	_       ErrorParser = new(Container)
)

// Container represents a HTML Form. The container can work with both HTTP Form and JSON requests
//
// swagger:model uiContainer
type Container struct {
	// Action should be used as the form action URL `<form action="{{ .Action }}" method="post">`.
	//
	// required: true
	Action string `json:"action" faker:"url"`

	// Method is the form method (e.g. POST)
	//
	// required: true
	Method string `json:"method" faker:"http_method"`

	// Nodes contains the form's nodes
	//
	// The form's nodes can be input fields, text, images, and other UI elements.
	//
	// required: true
	Nodes node.Nodes `json:"nodes"`

	// Messages contains all global form messages and errors.
	Messages text.Messages `json:"messages,omitempty"`
}
