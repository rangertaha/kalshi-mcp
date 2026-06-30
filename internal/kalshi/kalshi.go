// SPDX-License-Identifier: MIT

// Package kalshi holds the connection to the Kalshi trading API that the
// per-area tool packages (markets, …) share.
package kalshi

import (
	"fmt"

	"github.com/rangertaha/kalshi-mcp/internal/client"
)

// Clients bundles the REST clients needed to reach the Kalshi API.
type Clients struct {
	// API reaches the Kalshi trade API host.
	API *client.Client
}

// NewClients builds the Kalshi API client for the given base URL.
//
// Public market-data endpoints (used by the read-only markets toolset) need no
// authentication, so keyID/privateKey are optional. When both are supplied,
// requests are signed with the Kalshi RSA-PSS scheme, enabling portfolio and
// order endpoints.
func NewClients(baseURL, keyID, privateKey string, opts ...client.Option) (*Clients, error) {
	var auth client.Authorizer // nil => unauthenticated (public market data)
	if keyID != "" && privateKey != "" {
		signer, err := client.NewKalshiAuthorizer(keyID, privateKey)
		if err != nil {
			return nil, err
		}
		auth = signer
	}
	base := append([]client.Option{client.WithUserAgent("kalshi-mcp")}, opts...)

	api, err := client.New(baseURL, auth, base...)
	if err != nil {
		return nil, fmt.Errorf("creating kalshi client: %w", err)
	}
	return &Clients{API: api}, nil
}
