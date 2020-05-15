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

package stock

import (
	"context"
	"encoding/json"
	"io"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/onwsk8r/goiex/pkg/core/stock/price"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
)

// Price exposes methods for accessing IEX Cloud stock price data.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#stock-prices.
type Price struct {
	client iexcloud.Client
	path   string
}

// NewPrice creates a new Price with the given client
func NewPrice(client iexcloud.Client) *Price {
	log.Trace().Msg("instantiating new stock.Price")
	return &Price{
		client: client,
		path:   "stock",
	}
}

// HistoricalDaily returns adjusted and unadjusted data with daily resolution.
// The historical prices endpoint at https://iexcloud.io/docs/api/#historical-prices is capable
// of returning both daily and intraday data, which contain different fields.
// This method will return an error if the response cannot be decoded or if the request appears to
// be for intraday data.
func (s *Price) HistoricalDaily(ctx context.Context, ticker string,
	params *HistoricalPriceParams) ([]price.Historical, error) {
	uri := []string{s.path, ticker, "chart"}
	var prices []price.Historical

	err := s.client.Get(ctx, uri, params.Values(), func(r io.ReadCloser) error {
		if params.Range == "dynamic" {
			log.Debug().Msg("historical: creating anonymous struct to handle dynamic range")
			res := struct {
				Range string              `json:"range"`
				Data  *[]price.Historical `json:"data"`
			}{Data: &prices}
			return json.NewDecoder(r).Decode(&res)
		} else {
			return json.NewDecoder(r).Decode(&prices)
		}
	})
	return prices, err
}

// Intraday returns aggregated intraday prices in one minute buckets.
// https://iexcloud.io/docs/api/#intraday-prices
func (s *Price) Intraday(ctx context.Context, ticker string,
	params *HistoricalPriceParams) (prices []price.Intraday, err error) {
	uri := []string{s.path, ticker, "intraday-prices"}
	err = s.client.Get(ctx, uri, params.Values(), func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&prices)
	})
	return
}

// PreviousDay returns previous day adjusted price data for one stock.
// Use PreviousDayMarket() to fetch an array of data for all stocks.
// https://iexcloud.io/docs/api/#previous-day-price
func (s *Price) PreviousDay(ctx context.Context, ticker string) (p price.PreviousDay, err error) {
	uri := []string{s.path, ticker, "previous"}
	err = s.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&p)
	})
	return
}

// PreviousDayMarket returns previous day adjusted price data for all stocks.
// Use PreviousDay() to fetch data for a single stock.
// https://iexcloud.io/docs/api/#previous-day-price
func (s *Price) PreviousDayMarket(ctx context.Context) (prices []price.PreviousDay, err error) {
	uri := []string{s.path, "market", "previous"}
	err = s.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&prices)
	})
	return
}
