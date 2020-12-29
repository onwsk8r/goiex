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

package stock

import (
	"encoding/json"
	"fmt"
	"time"
)

// Split represents a data point from the basic splits endpoint.
// https://iexcloud.io/docs/api/#splits-basic
type Split struct {
	ExDate       time.Time `json:"exDate,omitempty"`
	DeclaredDate time.Time `json:"declaredDate,omitempty"`
	Ratio        float64   `json:"ratio,omitempty"`
	ToFactor     int       `json:"toFactor,omitempty"`
	FromFactor   int       `json:"fromFactor,omitempty"`
	Description  string    `json:"description,omitempty"`
	Symbol       string    `json:"symbol,omitempty"`
	ID           string    `json:"id,omitempty"`
	Source       string    `json:"source,omitempty"`
	Key          string    `json:"key,omitempty"`
	Subkey       string    `json:"subkey,omitempty"`
	Date         time.Time `json:"date,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (s *Split) UnmarshalJSON(data []byte) (err error) { // nolint:dupl
	type split Split
	type embedded struct {
		split
		ExDate       string `json:"exDate,omitempty"`
		DeclaredDate string `json:"declaredDate,omitempty"`
		Date         int64  `json:"date,omitempty"`
		Updated      int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*s = Split(tmp.split)
		// Ignore date parsing issues in case one or more dates are missing
		s.ExDate, _ = time.Parse("2006-01-02", tmp.ExDate)             // nolint:errcheck
		s.DeclaredDate, _ = time.Parse("2006-01-02", tmp.DeclaredDate) // nolint:errcheck
		if tmp.Date > 0 {
			s.Date = time.Unix(tmp.Date/1000, tmp.Date%1000*1e6) // nolint:gomnd
		}
		if tmp.Updated > 0 {
			s.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6) // nolint:gomnd
		}
	}
	return
}

// MarshalJSON satisfies the json.Marshaler interface.
// It undoes what UnmarshalJSON does.
func (s *Split) MarshalJSON() ([]byte, error) {
	type split Split
	type embedded struct {
		split
		ExDate       string `json:"exDate,omitempty"`
		DeclaredDate string `json:"declaredDate,omitempty"`
		Date         int64  `json:"date,omitempty"`
		Updated      int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	tmp.split = split(*s)
	tmp.DeclaredDate = s.DeclaredDate.Format("2006-01-02")
	tmp.ExDate = s.ExDate.Format("2006-01-02")
	tmp.Date = s.Date.UnixNano() / 1e6       // nolint:gomnd
	tmp.Updated = s.Updated.UnixNano() / 1e6 // nolint:gomnd
	return json.Marshal(tmp)
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
