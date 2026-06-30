// SPDX-License-Identifier: MIT

// Package markets exposes Kalshi public market data: listing markets and
// fetching a single market by ticker.
package markets

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rangertaha/kalshi-mcp/internal/kalshi"
)

// Name is the toolset name used for enable/disable filtering.
const Name = "markets"

// service wraps the Kalshi clients for market operations.
type service struct {
	c *kalshi.Clients
}

// Market is a Kalshi market, trimmed to the fields useful to an LLM. Prices are
// integer cents (0–100) representing the implied probability.
type Market struct {
	Ticker      string `json:"ticker"`
	EventTicker string `json:"event_ticker,omitempty"`
	Title       string `json:"title,omitempty"`
	Subtitle    string `json:"subtitle,omitempty"`
	Status      string `json:"status,omitempty"`
	YesBid      int    `json:"yes_bid"`
	YesAsk      int    `json:"yes_ask"`
	NoBid       int    `json:"no_bid"`
	NoAsk       int    `json:"no_ask"`
	LastPrice   int    `json:"last_price"`
	Volume      int    `json:"volume"`
	OpenTime    string `json:"open_time,omitempty"`
	CloseTime   string `json:"close_time,omitempty"`
}

// listResponse is the envelope returned by the markets list endpoint.
type listResponse struct {
	Markets []Market `json:"markets"`
	Cursor  string   `json:"cursor,omitempty"`
}

// getResponse is the envelope returned by the single-market endpoint.
type getResponse struct {
	Market Market `json:"market"`
}

// ListMarkets returns markets, optionally filtered by status and event, paged
// by limit/cursor.
func (s *service) ListMarkets(ctx context.Context, status, eventTicker, cursor string, limit int) ([]Market, error) {
	q := url.Values{}
	if status != "" {
		q.Set("status", status)
	}
	if eventTicker != "" {
		q.Set("event_ticker", eventTicker)
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var out listResponse
	if err := s.c.API.GetJSON(ctx, "/markets", q, &out); err != nil {
		return nil, err
	}
	return out.Markets, nil
}

// GetMarket returns a single market by its ticker.
func (s *service) GetMarket(ctx context.Context, ticker string) (*Market, error) {
	var out getResponse
	path := fmt.Sprintf("/markets/%s", url.PathEscape(ticker))
	if err := s.c.API.GetJSON(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out.Market, nil
}
