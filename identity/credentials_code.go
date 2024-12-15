// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identity

type CodeChannel string

const (
	CodeChannelEmail CodeChannel = AddressTypeEmail
	CodeChannelSMS   CodeChannel = AddressTypeSMS
)
