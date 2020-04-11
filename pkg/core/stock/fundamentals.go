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
	"strconv"

	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
	"github.com/rs/zerolog/log"
)

// Fundamentals exposes methods for accessing IEX Cloud fundamentals data.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#stock-fundamentals.
type Fundamentals struct {
	client iexcloud.Client
	path   string
}

// NewFundamentals creates a new Fundamentals with the given client.
func NewFundamentals(client iexcloud.Client) *Fundamentals {
	log.Trace().Msg("instantiating new stock.Fundamentals")
	return &Fundamentals{
		client: client,
		path:   "stock",
	}
}

// Dividends provides basic dividend data for US equities, ETFs, and Mutual Funds for the last 5 years.
// https://iexcloud.io/docs/api/#dividends-basic
func (f *Fundamentals) Dividends(ctx context.Context,
	symbol, rng string) (dividends []fundamental.Dividend, err error) {
	uri := []string{f.path, symbol, "dividends", rng}
	err = f.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&dividends)
	})
	return
}

// Earnings data for a given company including the actual EPS, consensus, and fiscal period.
// See https://iexcloud.io/docs/api/#earnings for more information. The "last" param specifies
// how many data points to return (up to 4). Set "annual" to true to return annual earnings
// instead of quarterly earnings.
func (f *Fundamentals) Earnings(ctx context.Context,
	ticker string, last int, annual bool) (earnings []fundamental.Earning, err error) {
	uri := []string{f.path, ticker, "earnings", strconv.Itoa(last)}
	v := url.Values{}
	if annual {
		v.Add("annual", "true")
	}
	err = f.client.Get(ctx, uri, v, func(r io.ReadCloser) error {
		res := struct {
			Symbol   string                 `json:"symbol"`
			Earnings *[]fundamental.Earning `json:"earnings"`
		}{Earnings: &earnings}
		return json.NewDecoder(r).Decode(&res)
	})
	return
}

// Splits returns up to five years of split history.
// https://iexcloud.io/docs/api/#splits-basic
func (f *Fundamentals) Splits(ctx context.Context, symbol, rng string) (splits []fundamental.Split, err error) {
	uri := []string{f.path, symbol, "splits", rng}
	err = f.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&splits)
	})
	return
}
