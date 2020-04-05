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

/* package api provides the reference implementation of the iexcloud.Client interface.
See the documentation for that package for more information on how a client might be
implemented in a manner that is consistent with the requirements and limitations of the
IEX API, and also for predefined variables and methods that may be useful in implementing
a clinet.*/
package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/onwsk8r/goiex/pkg/iexcloud"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context/ctxhttp"
)

// RequestsPerSecond limits the number of requests a Client instance can make in one second.
var RequestsPerSecond time.Duration = 50

// Client implements the iexcloud.Client interface.
// It provides a functional, rate-limited HTTP client that handles non-200
// HTTP responses. Requests use the ctxhttp library to be context-aware, and
// rate limiting is handled by a time.Ticker that ticks at a rate of
// `time.Second / RequestsPerSecond`.
type Client struct {
	client  *http.Client
	ticker  *time.Ticker
	domain  string
	version string
	token   string
}

// NewClient creates a new Client.
// The token is used directly, and if a sandbox token (ie one that starts with T) is passed,
// the client will use the sandbox domain. An empty token will cause the function to panic. The
// function will choose either the Base or Sandbox domain depending on whether the token is a
// sandbox token or not, respectively. If an empty version is passed, it will default to "stable".
// This function will log the token if trace logging is enabled.
func NewClient(token, version string) (*Client, error) {
	if version == "" {
		log.Info().Msg("client: received empty version: using 'stable'")
		version = iexcloud.APIVersionStable
	}
	log.Trace().Str("token", token).Str("version", version).Msg("client: creating new client")

	domain := iexcloud.APIDomainBase
	if token[0:4] == "Tpk_" || token[0:4] == "Tsk_" {
		log.Info().Msg("client: received sandbox token, using sandbox URL")
		domain = iexcloud.APIDomainSandbox
	}
	log.Debug().Str("domain", domain).Str("version", version).Msg("client: finalized params")

	return &Client{
		client:  http.DefaultClient,
		ticker:  time.NewTicker(time.Second / RequestsPerSecond),
		domain:  domain,
		version: version,
		token:   token,
	}, nil
}

// Get fulfills the iexcloud.Client interface.
// This method waits for either rate limiting or context.Done() before doing
// any work. The client's token is added to the query string parameters automatically,
// and the response body is automatically closed. If a non-200 repsonse status is received,
// this function will return an error without executing the passed function.
// This method will log the token if trace logging is enabled.
func (c *Client) Get(ctx context.Context, uri []string, params url.Values, f func(io.ReadCloser) error) error {
	log.Debug().Strs("uri", uri).Interface("qsp", params).Msg("client: initiating GET request")
	select {
	case <-ctx.Done():
		log.Debug().Msg("client: context is done prior to making request")
		return ctx.Err()
	case <-c.ticker.C:
		log.Trace().Msg("client: ticker has ticked")
	}

	params.Add("token", c.token)
	uri = append([]string{c.version}, uri...)
	url := fmt.Sprintf("%s/%s?%s", c.domain, path.Join(uri...), params.Encode())
	log.Trace().Str("url", url).Msg("client: performing GET request")

	res, err := ctxhttp.Get(ctx, c.client, url)
	if err != nil {
		log.Debug().Msg("client: received error from ctxhttp")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		log.Trace().Msg("client: calling passed function")
		return f(res.Body)
	}
	return fmt.Errorf("received HTTP status code %d (%s)", res.StatusCode, res.Status)
}
