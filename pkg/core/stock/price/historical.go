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
	"time"

	"github.com/rs/zerolog/log"
)

// Historical represents a data point from the Historical Prices endpoint.
// https://iexcloud.io/docs/api/#historical-prices
type Historical struct {
	Date           time.Time `json:"date"`
	Open           float64   `json:"open"`
	High           float64   `json:"high"`
	Low            float64   `json:"low"`
	Close          float64   `json:"close"`
	Volume         float64   `json:"volume"`
	UOpen          float64   `json:"uOpen"`
	UHigh          float64   `json:"uHigh"`
	ULow           float64   `json:"uLow"`
	UClose         float64   `json:"uClose"`
	UVolume        int       `json:"uVolume"`
	Change         float64   `json:"change"`
	ChangePercent  float64   `json:"changePercent"`
	Label          string    `json:"label"`
	ChangeOverTime float64   `json:"changeOverTime"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled or if the date parsing fails.
func (h *Historical) UnmarshalJSON(data []byte) (err error) {
	type historical Historical
	type embedded struct {
		historical
		Date string `json:"date"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*h = Historical(tmp.historical)
		h.Date, err = time.Parse("2006-01-02", tmp.Date)
		log.Debug().Str("original", tmp.Date).Time("parsed", h.Date).Msg("historical: parsed date")
	}
	return
}

// Validate satisfies the Validator interface.
// It will return an error if the Date, Close or UClose fields are equal to their zero value.
func (h *Historical) Validate() error {
	switch {
	case h.Date.IsZero():
		return fmt.Errorf("missing date")
	case h.Close == 0:
		return fmt.Errorf("close is zero")
	case h.UClose == 0:
		return fmt.Errorf("unadjusted close is zero")
	}
	return nil
}
