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

package core

import (
	"context"
	"encoding/json"
	"io"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/onwsk8r/goiex/pkg/core/market"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
)

// Market exposes methods for accessing IEX Cloud market information.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#market-info.
type Market struct {
	client iexcloud.Client
}

// NewMarket creates a new Market with the given client
func NewMarket(client iexcloud.Client) *Market {
	log.Trace().Msg("instantiating new core.Market")
	return &Market{
		client: client,
	}
}

// UpcomingDividends fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingDividends(ctx context.Context,
	symbol string) (dividends []market.UpcomingDividend, err error) {
	uri := []string{"stock", symbol, "upcoming-dividends"}
	err = m.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&dividends)
	})
	return
}

// UpcomingEarnings fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingEarnings(ctx context.Context, symbol string) (earnings []market.UpcomingEarning, err error) {
	uri := []string{"stock", symbol, "upcoming-earnings"}
	err = m.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&earnings)
	})
	return
}

// UpcomingSplits fetches a list of upcoming earnings from the upcoming events endpoint.
// Pass "market" as the symbol parameter to return data for all symbols.
// https://iexcloud.io/docs/api/#upcoming-events
func (m *Market) UpcomingSplits(ctx context.Context, symbol string) (splits []market.UpcomingSplit, err error) {
	uri := []string{"stock", symbol, "upcoming-splits"}
	err = m.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&splits)
	})
	return
}
