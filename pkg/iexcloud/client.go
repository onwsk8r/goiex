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

/* package iexcloud provides primitives for working with the IEX Cloud API. Nothing in this package
provides de facto functionality, but the contents can be used as standardized building blocks for
implementing API clients. A reference implementation of the contents of this package can be found in
the neighboring internal/api package.*/
package iexcloud

import (
	"context"
	"io"
	"net/url"
)

// These variables correspond to the domains the IEX Cloud API uses to serve content.
var (
	APIDomainBase       string = "https://cloud.iexapis.com"
	APIDomainSandbox    string = "https://sandbox.iexapis.com"
	APIDomainSandboxSSE string = "https://sandbox-sse.iexapis.com"
)

// These variables correspond to the available versions (ie the first path segment) of the IEX Cloud API.
var (
	APIVersion1      string = "v1"
	APIVersionBeta   string = "beta"
	APIVersionStable string = "stable"
	APIVersionLatest string = "latest"
)

// Client provides a consistent interface for HTTP clients.
// Each of the types in the adjacent packages that make HTTP requests will consume
// this interface and use it to make requests. An instance of Client corresponds to a particular
// API token and version. An implementation of Client should, besides fulfilling the interface,
// - Keep the token abstracted away from code that makes requests
// - Handle rate limiting and non-200 response statuses
// - Handle context cancellation
// - Be thread-safe
// As this is not a generic http.Client - it is specific to this use-case - the parameters and
// expected values are use-case specific as well. Since the IEX API only supports the HTTP GET
// method, the interface only defines a Get method. The parameters to this function are, in order:
// - A context that, if Done(), will cause the function to cancel in-flight requests and return an error
// - An array of path segments, such as []string{"stock", "AAPL", "price"}, to be appended to the version
// - A url.Values{} containing any query string parameters. The Client will Add() the token.
// - A function that consumes the http.Response.Body, which is only called with a HTTP 200 return status.
type Client interface {
	Get(ctx context.Context, segments []string, params url.Values, f func(io.ReadCloser) error) error
}
