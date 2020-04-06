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

package price

import (
	"encoding/json"
	"fmt"
)

// PreviousDay represents a data point from the Previous Day Prices endpoint.
// https://iexcloud.io/docs/api/#previous-day-price
type PreviousDay struct {
	Historical
	Symbol string `json:"symbol"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function works around the inherited UnmarshalJson() from the embedded Historical
func (p *PreviousDay) UnmarshalJSON(data []byte) (err error) {
	// First get the Historical part
	var h Historical
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}
	p.Historical = h

	// And then that left over ticker symbol
	type ticker struct {
		Symbol string `json:"symbol"`
	}
	tmp := new(ticker)
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}
	p.Symbol = tmp.Symbol
	return nil
}

// Validate satisfies the Validator interface.
// It will return an error if the Symbol is zero-valued or if the embedded Historical is invalid.
func (p *PreviousDay) Validate() error {
	if p.Symbol == "" {
		return fmt.Errorf("missing symbol")
	}
	return p.Historical.Validate()
}
