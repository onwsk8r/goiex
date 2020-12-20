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

package option

import (
	"time"
)

// OptionOrig represents the original options data format.
// This schema was deprecated on Dec 1, 2020, and remains for compatibility
// because the new format makes breaking changes.
type OptionOrig struct {
	Symbol         string    `json:"symbol"`
	ID             string    `json:"id"`
	ExpirationDate time.Time `json:"expirationDate"`
	ContractSize   int       `json:"contractSize"`
	StrikePrice    float64   `json:"strikePrice"`
	ClosingPrice   float64   `json:"closingPrice"`
	Side           string    `json:"side"`
	Type           string    `json:"type"`
	Volume         int       `json:"volume"`
	OpenInterest   int       `json:"openInterest"`
	Bid            float64   `json:"bid"`
	Ask            float64   `json:"ask"`
	LastUpdated    time.Time `json:"lastUpdated"`
	IsAdjusted     bool      `json:"isAdjusted"`
}
