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
	"encoding/json"
	"fmt"

	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
)

// UpcomingSplit represents a data point from the upcoming events endpoint.
// The data structure is similar to the basic splits endpoint with an added symbol field. See
// https://iexcloud.io/docs/api/#splits-basic and
// https://iexcloud.io/docs/api/#upcoming-events for more information.
type UpcomingSplit struct {
	fundamental.Split
	Symbol string `json:"symbol"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function works around the inherited UnmarshalJson() from the embedded Split
func (u *UpcomingSplit) UnmarshalJSON(data []byte) (err error) {
	// First get the Split part
	var s fundamental.Split
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	u.Split = s

	// And then that left over ticker symbol
	type ticker struct {
		Symbol string `json:"symbol"`
	}
	tmp := new(ticker)
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}
	u.Symbol = tmp.Symbol
	return nil
}

// Validate satisfies the Validator interface.
// It will return an error if the DeclaredDate or ExDate are zero, or the Symbol is missing
func (u *UpcomingSplit) Validate() (err error) {
	switch {
	case u.Symbol == "":
		err = fmt.Errorf("symbol is missing")
	case u.DeclaredDate.IsZero():
		err = fmt.Errorf("declared date is missing")
	case u.ExDate.IsZero():
		err = fmt.Errorf("ex date is missing")
	}
	return
}
