// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package node

func NewCSRFNode(token string) *Node {
	return &Node{
		Type:  Input,
		Group: DefaultGroup,
	}
}
