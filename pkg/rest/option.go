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
	"github.com/onwsk8r/goiex/pkg/core/option"
)

// Options exposes methods for calling options endpoints.
// https://iexcloud.io/docs/api/#options
type Options struct {
	client *resty.Client
}

// NewOptions creates a new Options with the given client
// https://iexcloud.io/docs/api/#options
func NewOptions(client *resty.Client) *Options {
	return &Options{client: client}
}

// Expiration returns a slice of expiration dates for the given symbol.
// This endpoint is defined as part of the End of Day Options endpoint, and is
// the result of a call to that endpoint that does not specify an expiration.
// https://iexcloud.io/docs/api/#end-of-day-options
func (o *Options) Expiration(ctx context.Context, symbol string) (res []string, err error) {
	var params = make(map[string]string)
	params["symbol"] = symbol

	_, err = o.client.R().SetContext(ctx).SetPathParams(params).SetResult(&res).
		Get("/{version}/stock/{symbol}/options")
	return
}

// EndOfDay returns a slice of option.Options for the given symbol and expiration.
// The expiration is expected to be specified as YYYYMM or YYYYMMDD: the format returned
// from this same endpoint when an expiration is not specified (see Expiration()) or
// the analogous Reference endpoint. If a side ("call" or "put") is passed as an
// optional fourth parameter, the method will return only options for that side.
// https://iexcloud.io/docs/api/#end-of-day-options
func (o *Options) EndOfDay(ctx context.Context, symbol, expiration string,
	side ...string) (res []option.Option, err error) {
	var params = make(map[string]string)
	params["symbol"] = symbol
	params["expiration"] = expiration
	params["side"] = ""
	if len(side) > 0 {
		params["side"] = side[0]
	}
	_, err = o.client.R().SetContext(ctx).SetPathParams(params).SetResult(&res).
		Get("/{version}/stock/{symbol}/options/{expiration}/{side}")
	return
}
