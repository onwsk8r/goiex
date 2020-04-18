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

package market

import (
	"fmt"

	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
)

// UpcomingDividend represents a data point from the upcoming events endpoint.
// The data structure is similar to the basic dividends endpoint with an added symbol field. See
// https://iexcloud.io/docs/api/#dividends-basic and
// https://iexcloud.io/docs/api/#upcoming-events for more information.
type UpcomingDividend struct {
	fundamental.Dividend
}

// Validate satisfies the Validator interface.
// It will return an error if the Symbol, declared date or ex-date are zero-valued. It will not
// validate the emdedded dividend because required values - such as the amount - may be missing.
func (u *UpcomingDividend) Validate() (err error) {
	switch {
	case u.Symbol == "":
		err = fmt.Errorf("symbol is missing")
	case u.ExDate.IsZero():
		err = fmt.Errorf("ex date is missing")
	case u.DeclaredDate.IsZero():
		err = fmt.Errorf("declared date is missing")
	}
	return err
}
