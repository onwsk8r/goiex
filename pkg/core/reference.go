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

	"github.com/onwsk8r/goiex/pkg/core/reference"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
	"github.com/rs/zerolog/log"
)

// Reference exposes methods for accessing IEX Cloud reference data.
// A list of endpoints can be found at https://iexcloud.io/docs/api/#reference-data.
type Reference struct {
	client iexcloud.Client
}

// NewReference creates a new Reference with the given client
func NewReference(client iexcloud.Client) *Reference {
	log.Trace().Msg("instantiating new core.Reference")
	return &Reference{client}
}

// Symbols fetches a list of symbols that are supported for making API calls.
// https://iexcloud.io/docs/api/#symbols
func (r *Reference) Symbols(ctx context.Context) (symbols []reference.Symbol, err error) {
	uri := []string{"ref-data", "symbols"}
	log.Trace().Strs("uri", uri).Msg("reference: calling client")
	err = r.client.Get(ctx, uri, url.Values{}, func(r io.ReadCloser) error {
		return json.NewDecoder(r).Decode(&symbols)
	})
	return
}
