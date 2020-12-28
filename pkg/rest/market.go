// goiex: Golang interface to IEX Cloud API
// Copyright (C) 2019 Brian Hazeltine

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rest

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/onwsk8r/goiex/pkg/core/market"
	"github.com/onwsk8r/goiex/pkg/core/stock"
)

// Market exposes methods for accessing IEX Cloud market information.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#market-info.
type Market struct {
	client *resty.Client
}

// NewMarket creates a new Market with the given client
func NewMarket(client *resty.Client) *Market {
	return &Market{
		client: client,
	}
}

// UpcomingDividends fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingDividends(ctx context.Context,
	symbol string) (dividends []market.UpcomingDividend, err error) {
	var params = map[string]string{"symbol": symbol}
	_, err = m.client.R().SetContext(ctx).SetPathParams(params).SetResult(&dividends).
		Get("/{version}/stock/{symbol}/upcoming-dividends")
	return
}

// UpcomingEarnings fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingEarnings(ctx context.Context, symbol string) (earnings []market.UpcomingEarning, err error) {
	var params = map[string]string{"symbol": symbol}
	_, err = m.client.R().SetContext(ctx).SetPathParams(params).SetResult(&earnings).
		Get("/{version}/stock/{symbol}/upcoming-earnings")
	return
}

// UpcomingSplits fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingSplits(ctx context.Context, symbol string) (splits []stock.Split, err error) {
	var params = map[string]string{"symbol": symbol}
	_, err = m.client.R().SetContext(ctx).SetPathParams(params).SetResult(&splits).
		Get("/{version}/stock/{symbol}/upcoming-splits")
	return
}
