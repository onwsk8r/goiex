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
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Option represents a data point from the Options EOD endpoint.
// https://iexcloud.io/docs/api/#end-of-day-options
type Option struct {
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

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the expirationDate and lastUpdated fields, which are specified as
// "YYYYMMDD" and "YYYY-MM-DD", respectively, into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the lastUpdated parsing fails.
func (o *Option) UnmarshalJSON(data []byte) (err error) {
	type option Option
	type embedded struct {
		option
		ExpirationDate string `json:"expirationDate"`
		LastUpdated    string `json:"lastUpdated"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*o = Option(tmp.option)
		o.ExpirationDate, _ = time.Parse("20060102", tmp.ExpirationDate) // nolint: errcheck
		o.LastUpdated, err = time.Parse("2006-01-02", tmp.LastUpdated)
		log.Debug().
			Dict("expiration_date", zerolog.Dict().Str("original", tmp.ExpirationDate).Time("parsed", o.ExpirationDate)). // nolint: lll
			Dict("expiration_date", zerolog.Dict().Str("original", tmp.LastUpdated).Time("parsed", o.LastUpdated)).
			Msg("option: parsed dates")
	}
	return
}

// Validate satisfies the Validator interface.
// It will return an error if the Date, Close or UClose fields are equal to their zero value.
func (o *Option) Validate() error {
	switch {
	case o.Symbol == "":
		return fmt.Errorf("missing symbol")
	case o.ID == "":
		return fmt.Errorf("missing id")
	case o.ExpirationDate.IsZero():
		return fmt.Errorf("missing expiration date")
	case o.StrikePrice == 0:
		return fmt.Errorf("strike price is zero")
	case o.ClosingPrice == 0:
		return fmt.Errorf("closing price is zero")
	}
	return nil
}
