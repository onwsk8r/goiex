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

package fundamental

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Split represents a data point from the basic splits endpoint.
// https://iexcloud.io/docs/api/#splits-basic
type Split struct {
	ExDate       time.Time `json:"exDate"`
	DeclaredDate time.Time `json:"declaredDate"`
	Ratio        float64   `json:"ratio"`
	ToFactor     int       `json:"toFactor"`
	FromFactor   int       `json:"fromFactor"`
	Description  string    `json:"description"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (s *Split) UnmarshalJSON(data []byte) (err error) {
	type split Split
	type embedded struct {
		split
		ExDate       string `json:"exDate"`
		DeclaredDate string `json:"declaredDate"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*s = Split(tmp.split)
		// Ignore date parsing issues in case one or more dates are missing
		s.ExDate, _ = time.Parse("2006-01-02", tmp.ExDate)             // nolint: errcheck
		s.DeclaredDate, _ = time.Parse("2006-01-02", tmp.DeclaredDate) // nolint: errcheck
		log.Debug().
			Dict("ex_date", zerolog.Dict().Str("original", tmp.ExDate).Time("parsed", s.ExDate)).
			Dict("declared_date", zerolog.Dict().Str("original", tmp.DeclaredDate).Time("parsed", s.DeclaredDate)).
			Msg("split: parsed date")
	}
	return
}

// Validate satisfies the Validator interface.
// It will return an error if the ExDate is equal to its zero value, or if the ToFactor or FromFactor are not positive.
func (s *Split) Validate() error {
	switch {
	case s.ToFactor <= 0:
		return fmt.Errorf("to factor is not positive")
	case s.FromFactor <= 0:
		return fmt.Errorf("from factor is not positive")
	case s.ExDate.IsZero():
		return fmt.Errorf("ex date is missing")
	}
	return nil
}
