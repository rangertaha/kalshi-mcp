// SPDX-License-Identifier: MIT

package markets

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/kalshi-mcp/internal/kalshi"
	"github.com/rangertaha/kalshi-mcp/internal/server"
)

// Register adds the markets toolset to the server.
func Register(s *server.Server, c *kalshi.Clients) {
	s.NoteToolset(Name)
	svc := &service{c: c}

	server.Register(s, server.ToolDef{
		Name:        "markets_list",
		Title:       "List markets",
		Description: "List Kalshi markets, optionally filtered by status (open, closed, settled) or event ticker. Prices are cents (0-100) implying probability.",
	}, svc.list)

	server.Register(s, server.ToolDef{
		Name:        "markets_get",
		Title:       "Get market",
		Description: "Get a single Kalshi market by its ticker.",
	}, svc.get)
}

// --- Tool input types (schemas are inferred from these structs) ---

// ListInput filters and pages the markets list.
type ListInput struct {
	Status      string `json:"status,omitempty" jsonschema:"filter by status: open, closed, or settled (optional)"`
	EventTicker string `json:"eventTicker,omitempty" jsonschema:"filter to markets under this event ticker (optional)"`
	Cursor      string `json:"cursor,omitempty" jsonschema:"pagination cursor from a previous response (optional)"`
	Limit       int    `json:"limit,omitempty" jsonschema:"maximum number of markets to return, 1-1000 (optional)"`
}

// GetInput identifies a single market.
type GetInput struct {
	Ticker string `json:"ticker" jsonschema:"market ticker, e.g. KXPRES-24"`
}

// --- Tool handlers ---

func (s *service) list(ctx context.Context, _ *mcp.CallToolRequest, in ListInput) (*mcp.CallToolResult, server.ListResult[Market], error) {
	out, err := s.ListMarkets(ctx, in.Status, in.EventTicker, in.Cursor, in.Limit)
	return nil, server.List(out), err
}

func (s *service) get(ctx context.Context, _ *mcp.CallToolRequest, in GetInput) (*mcp.CallToolResult, *Market, error) {
	out, err := s.GetMarket(ctx, in.Ticker)
	return nil, out, err
}
