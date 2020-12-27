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

	"github.com/rs/zerolog/log"
)

// Historical represents a data point from the Historical Prices endpoint.
// https://iexcloud.io/docs/api/#historical-prices
type Historical struct {
	Close                float64   `json:"close,omitempty"`
	High                 float64   `json:"high,omitempty"`
	Low                  float64   `json:"low,omitempty"`
	Open                 float64   `json:"open,omitempty"`
	Symbol               string    `json:"symbol,omitempty"`
	Volume               float64   `json:"volume,omitempty"`
	ID                   string    `json:"id,omitempty"`
	Source               string    `json:"source,omitempty"`
	Key                  string    `json:"key,omitempty"`
	Subkey               string    `json:"subkey,omitempty"`
	Date                 time.Time `json:"date,omitempty"`
	Updated              time.Time `json:"updated,omitempty"`
	ChangeOverTime       float64   `json:"changeOverTime,omitempty"`
	MarketChangeOverTime float64   `json:"marketChangeOverTime,omitempty"`
	UOpen                float64   `json:"uOpen,omitempty"`
	UHigh                float64   `json:"uHigh,omitempty"`
	ULow                 float64   `json:"uLow,omitempty"`
	UClose               float64   `json:"uClose,omitempty"`
	UVolume              uint      `json:"uVolume,omitempty"`
	FOpen                float64   `json:"fOpen,omitempty"`
	FHigh                float64   `json:"fHigh,omitempty"`
	FLow                 float64   `json:"fLow,omitempty"`
	FClose               float64   `json:"fClose,omitempty"`
	FVolume              float64   `json:"fVolume,omitempty"`
	Label                string    `json:"label,omitempty"`
	Change               float64   `json:"change,omitempty"`
	ChangePercent        float64   `json:"changePercent,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (h *Historical) UnmarshalJSON(data []byte) (err error) {
	type historical Historical
	type embedded struct {
		historical
		Date    int64 `json:"date,omitempty"`
		Updated int64 `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*h = Historical(tmp.historical)
		// Ignore date parsing issues, which will happen especially on PreviousDayMarket
		// calls where some array objects may be empty (ie "{}")
		h.Date = time.Unix(tmp.Date/1000, tmp.Date%1000*1e6)          // nolint:gomnd
		h.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6) // nolint:gomnd
		log.Debug().Interface("original", tmp).Interface("final", h).Msg("historical: parsed date")
	}
	return
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
