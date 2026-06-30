// SPDX-License-Identifier: MIT

package kalshi

import (
	"context"
	"net/url"
)

// Check verifies connectivity by requesting a single market. It returns the
// number of markets in the returned page.
func Check(ctx context.Context, c *Clients) (int, error) {
	q := url.Values{}
	q.Set("limit", "1")
	var out struct {
		Markets []struct {
			Ticker string `json:"ticker"`
		} `json:"markets"`
	}
	if err := c.API.GetJSON(ctx, "/markets", q, &out); err != nil {
		return 0, err
	}
	return len(out.Markets), nil
}
