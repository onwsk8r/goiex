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

package reference

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// Symbol represents one datum of that returned by the ref-data/symbols endpoint.
type Symbol struct {
	Symbol   string    `json:"symbol"`
	Name     string    `json:"name"`
	Exchange string    `json:"exchange"`
	IEXID    string    `json:"iexid"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"-"`
	Type     string    `json:"type"`
	Enabled  bool      `json:"isEnabled"`
	Region   string    `json:"region"`
	FIGI     string    `json:"figi"`
	CIK      string    `json:"cik"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// into a time.Time by using time.Parse().
func (s *Symbol) UnmarshalJSON(data []byte) (err error) {
	type symbol Symbol
	type embedded struct {
		symbol
		Date string `json:"date"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*s = Symbol(tmp.symbol)
		s.Date, err = time.Parse("2006-01-02", tmp.Date)
		log.Debug().Str("original", tmp.Date).Time("parsed", s.Date).Msg("symbol: parsed date")
	}
	return
}

// MarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// into a time.Time by using time.Parse().
func (s *Symbol) MarshalJSON() ([]byte, error) {
	type symbol Symbol
	type embedded struct {
		symbol
		Date string `json:"date"`
	}
	tmp := new(embedded)
	tmp.symbol = symbol(*s)
	tmp.Date = s.Date.Format("2006-01-02")
	return json.Marshal(tmp)
}

// Validate satisfies the Validator interface.
// It will return an error if the Symbol, IEXID, or Date fields are equal to their zero value.
func (s *Symbol) Validate() error {
	switch {
	case s.Symbol == "":
		return fmt.Errorf("missing symbol")
	case s.IEXID == "":
		return fmt.Errorf("missing IEX ID")
	case s.Date.IsZero():
		return fmt.Errorf("missing date")
	}
	return nil
}
