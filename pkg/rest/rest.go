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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// These variables correspond to the domains the IEX Cloud API uses to serve content.
var (
	APIDomainBase    string = "https://cloud.iexapis.com"
	APIDomainSandbox string = "https://sandbox.iexapis.com"
)

// These variables correspond to the available versions (ie the first path segment) of the IEX Cloud API.
var (
	APIVersion1      string = "v1"
	APIVersionBeta   string = "beta"
	APIVersionStable string = "stable"
	APIVersionLatest string = "latest"
)

// RequestsPerSecond is used by NewClient to implement rate limiting
var RequestsPerSecond = 50

// MaxRetries is used by NewClient to limit the number of retry attempts
var MaxRetries = 5

var ticker *time.Ticker = nil

// NewClient creates a new go-resty client with some helpful configuration.
// - The token is set as a query string parameter, and the HostURL (ie domain) is
// initialized to the regular or sandbox domain accordingly. A "version" path parameter
// is set to "v1" by default. Both are set from package variables.
// - Passing a non-nil Logger will set it as the go-resty logger.
// - It implements rate limiting with an unexported time.Ticker, set at the package
// level, that is reused for each client via OnBeforeRequest. This RPS limiter
// is not related to exponential backoff for retries.
// - The MaxRetries package variable sets the retry count, and a RetryConditionFunc
// returns true if Response.IsError() with a status code > 404 and not 413 or 451.
// - If Response.IsError(), an error will be returned that matches the format
// "invalid response: <code> <response-body>" (eg "invalid response: 404 Unknown Symbol")
func NewClient(token string, logger *zerolog.Logger) *resty.Client {
	if ticker == nil {
		ticker = time.NewTicker(time.Second / time.Duration(RequestsPerSecond))
	}

	c := resty.New().
		AddRetryCondition(checkRetry).
		OnBeforeRequest(requestLimiter).
		OnAfterResponse(checkResponse).
		SetPathParams(map[string]string{"version": APIVersion1}).
		SetQueryParam("token", token).
		SetRetryCount(MaxRetries)

	if logger != nil {
		c.SetLogger(zl{l: logger})
	}

	// Only sandbox tokens start with "T" (Tsk_, Tpk_ vs sk_, pk_)
	if strings.Index(token, "T") == 0 {
		return c.SetHostURL(APIDomainSandbox)
	}
	return c.SetHostURL(APIDomainBase)
}

// requestLimiter uses a time.Ticker to limit requests per second
func requestLimiter(_ *resty.Client, req *resty.Request) error {
	select {
	case <-ticker.C:
	case <-req.Context().Done():
		return req.Context().Err()
	}
	return nil
}

// checkResponse sets an error for HTTP status codes >=400 (ie resp.IsError())
func checkResponse(c *resty.Client, resp *resty.Response) error {
	if resp.IsError() {
		return fmt.Errorf("invalid response: %s %s", resp.Status(), resp)
	}
	return nil
}

// checkRetry returns true for non-nil erros, except for canceled contexts and 400-404,413,451
func checkRetry(r *resty.Response, err error) bool {
	if r != nil && r.IsError() && (r.StatusCode() == http.StatusBadRequest ||
		r.StatusCode() == http.StatusUnauthorized ||
		r.StatusCode() == http.StatusPaymentRequired ||
		r.StatusCode() == http.StatusForbidden ||
		r.StatusCode() == http.StatusNotFound ||
		r.StatusCode() == http.StatusRequestEntityTooLarge ||
		r.StatusCode() == http.StatusUnavailableForLegalReasons) {
		return false
	}
	return err != nil && err.Error() != "context canceled"
}

type zl struct {
	l *zerolog.Logger
}

func (z zl) Errorf(format string, v ...interface{}) { z.l.Error().Msgf(format, v...) }
func (z zl) Warnf(format string, v ...interface{})  { z.l.Warn().Msgf(format, v...) }
func (z zl) Debugf(format string, v ...interface{}) { z.l.Debug().Msgf(format, v...) }
