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

	"github.com/onwsk8r/goiex/pkg/core/option"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
	"github.com/rs/zerolog/log"
)

// Option exposes methods for accessing IEX Cloud EOD Option data.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#options.
type Option struct {
	client iexcloud.Client
}

// NewOption creates a new Option with the given client
func NewOption(client iexcloud.Client) *Option {
	log.Trace().Msg("instantiating new core.Option")
	return &Option{client}
}

// Options fetches a list of options for a given symbol and expiration date.
// For consistency with the reference endpoint that provides a list of expiration months,
// the date is expected to be in the "YYYYMM" format returned by that endpoint.
// https://iexcloud.io/docs/api/#end-of-day-options
func (o *Option) Options(ctx context.Context, symbol, date string) (options []option.Option, err error) {
	uri := []string{"stock", symbol, "options", date}
	log.Trace().Strs("uri", uri).Msg("option: calling client")
	err = o.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&options)
	})
	return
}
