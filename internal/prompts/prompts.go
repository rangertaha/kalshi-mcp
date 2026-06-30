// SPDX-License-Identifier: MIT

// Package prompts registers MCP prompts: user-invoked, parameterized templates
// that clients surface as slash commands. Each prompt encodes a multi-step
// workflow by guiding the model to call the right tools in order.
package prompts

import (
	"fmt"

	"github.com/rangertaha/kalshi-mcp/internal/server"
)

// Register adds the built-in workflow prompts to the server.
func Register(s *server.Server) {
	s.AddPrompt(
		"market_odds",
		"Read the implied odds for a Kalshi market and explain what the prices mean.",
		[]server.PromptArg{
			{Name: "ticker", Description: "market ticker", Required: true},
		},
		func(a map[string]string) string {
			return fmt.Sprintf(`Explain the odds for Kalshi market "%s".

Steps:
1. Call markets_get (ticker="%s") to load the market.
2. Kalshi prices are cents (0-100) that imply probability. Report the current
   yes bid/ask and what the midpoint implies as a probability of the event.
3. Note the market status, volume, and close time, and flag if the market is
   already closed or settled.`,
				a["ticker"], a["ticker"])
		},
	)
}
