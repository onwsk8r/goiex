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

// Historical represents a data point from the Historical Prices endpoint.
// https://iexcloud.io/docs/api/#historical-prices
type Historical struct {
	Close                float64   `json:"close,omitempty" gorm:"type:double precision"`
	High                 float64   `json:"high,omitempty" gorm:"type:double precision"`
	Low                  float64   `json:"low,omitempty" gorm:"type:double precision"`
	Open                 float64   `json:"open,omitempty" gorm:"type:double precision"`
	Symbol               string    `json:"symbol,omitempty" gorm:"primaryKey;type:character varying"`
	Volume               *float64  `json:"volume,omitempty" gorm:"type:double precision"`
	ID                   string    `json:"id,omitempty" gorm:"-"`
	Source               string    `json:"source,omitempty" gorm:"-"`
	Key                  string    `json:"key,omitempty" gorm:"-"`
	Subkey               string    `json:"subkey,omitempty" gorm:"-"`
	Date                 time.Time `json:"-" gorm:"primaryKey"`
	Updated              time.Time `json:"-"`
	ChangeOverTime       *float64  `json:"changeOverTime,omitempty" gorm:"type:double precision"`
	MarketChangeOverTime *float64  `json:"marketChangeOverTime,omitempty" gorm:"type:double precision"`
	UOpen                float64   `json:"uOpen,omitempty" gorm:"type:double precision"`
	UHigh                float64   `json:"uHigh,omitempty" gorm:"type:double precision"`
	ULow                 float64   `json:"uLow,omitempty" gorm:"type:double precision"`
	UClose               float64   `json:"uClose,omitempty" gorm:"type:double precision"`
	UVolume              *float64  `json:"uVolume,omitempty" gorm:"type:double precision"`
	FOpen                float64   `json:"fOpen,omitempty" gorm:"type:double precision"`
	FHigh                float64   `json:"fHigh,omitempty" gorm:"type:double precision"`
	FLow                 float64   `json:"fLow,omitempty" gorm:"type:double precision"`
	FClose               float64   `json:"fClose,omitempty" gorm:"type:double precision"`
	FVolume              *float64  `json:"fVolume,omitempty" gorm:"type:double precision"`
	Label                string    `json:"label,omitempty" gorm:"type:character varying"`
	Change               *float64  `json:"change,omitempty" gorm:"type:double precision"`
	ChangePercent        *float64  `json:"changePercent,omitempty" gorm:"type:double precision"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (h *Historical) UnmarshalJSON(data []byte) (err error) {
	type historical Historical
	type embedded struct {
		historical
		Date    string `json:"date,omitempty"`
		Updated int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*h = Historical(tmp.historical)
		h.Date, _ = time.Parse("2006-01-02", tmp.Date)                // nolint:errcheck
		h.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6) // nolint:gomnd
	}
	return
}

// MarshalJSON satisfies the json.Marshaler interface.
// It undoes what UnmarshalJSON does.
func (h *Historical) MarshalJSON() ([]byte, error) {
	type historical Historical
	type embedded struct {
		historical
		Date    string `json:"date,omitempty"`
		Updated int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	tmp.historical = historical(*h)
	tmp.Date = h.Date.Format("2006-01-02")
	tmp.Updated = h.Updated.UnixNano() / 1e6 // nolint:gomnd
	return json.Marshal(tmp)
}

// Validate satisfies the Validator interface.
// It will return an error if the Date, Close or UClose fields are equal to their zero value.
func (h *Historical) Validate() error {
	switch {
	case h.Symbol == "":
		return fmt.Errorf("missing symbol")
	case h.Date.IsZero():
		return fmt.Errorf("missing date")
	case h.Close == 0:
		return fmt.Errorf("close is zero")
	case h.UClose == 0:
		return fmt.Errorf("unadjusted close is zero")
	}
	return nil
}
