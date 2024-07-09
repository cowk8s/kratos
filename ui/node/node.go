// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package node

// swagger:enum UiNodeType
type UiNodeType string

const (
	Text   UiNodeType = "text"
	Input  UiNodeType = "input"
	Image  UiNodeType = "img"
	Anchor UiNodeType = "a"
	Script UiNodeType = "script"
)

type Nodes []Node

// Node represents a flow's nodes
//
// Nodes are represented as HTML elements or their native UI equivalents. For example,
// a node can be an `<img>` tag, or an `<input element>` but also `some plain text`.
//
// swagger:model uiNode
type Node struct {
	// The node's type
	//
	// required: true
	Type UiNodeType `json:"type" faker:"-"`
}
