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
	"github.com/onwsk8r/goiex/pkg/core/reference"
)

// Reference exposes methods for calling Reference endoints.
// https://iexcloud.io/docs/api/#reference-data
type Reference struct {
	client *resty.Client
}

// NewReference creates a new Reference with the given client
// https://iexcloud.io/docs/api/#reference-data
func NewReference(client *resty.Client) *Reference {
	return &Reference{client: client}
}

// Symbols fetches a list of symbols that are supported for making API calls.
// The list seems to be exclusive to US equities.
// https://iexcloud.io/docs/api/#symbols
func (r *Reference) Symbols(ctx context.Context) (symbols []reference.Symbol, err error) {
	_, err = r.client.R().SetContext(ctx).SetResult(&symbols).Get("/{version}/ref-data/symbols")
	return
}

// OptionsSymbols fetches a list of options symbols/dates that are supported for making API calls.
// This call returns an object keyed by symbol with the value of each symbol being an array of available contract dates.
// https://iexcloud.io/docs/api/#options-symbols
func (r *Reference) OptionsSymbols(ctx context.Context) (symbols reference.OptionSymbol, err error) {
	_, err = r.client.R().SetContext(ctx).SetResult(&symbols).Get("/{version}/ref-data/options/symbols")
	return
}
