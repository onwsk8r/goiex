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

// Intraday represents a data point from the Intraday Prices endpoint.
// https://iexcloud.io/docs/api/#intraday-prices
type Intraday struct {
	Symbol               string    `json:"symbol,omitempty" gorm:"primaryKey;type:character varying"`
	Date                 time.Time `json:"date,omitempty" gorm:"primaryKey"`
	Minute               string    `json:"minute,omitempty" gorm:"type:character varying"`
	Label                string    `json:"label,omitempty" gorm:"type:character varying"`
	MarketOpen           float64   `json:"marketOpen,omitempty" gorm:"type:double precision"`
	MarketHigh           float64   `json:"marketHigh,omitempty" gorm:"type:double precision"`
	MarketLow            float64   `json:"marketLow,omitempty" gorm:"type:double precision"`
	MarketClose          float64   `json:"marketClose,omitempty" gorm:"type:double precision"`
	MarketVolume         int       `json:"marketVolume,omitempty"`
	MarketAverage        float64   `json:"marketAverage,omitempty" gorm:"type:double precision"`
	MarketNotional       float64   `json:"marketNotional,omitempty" gorm:"type:double precision"`
	MarketNumberOfTrades int       `json:"marketNumberOfTrades,omitempty"`
	MarketChangeOverTime float64   `json:"marketChangeOverTime,omitempty" gorm:"type:double precision"`
	Open                 float64   `json:"open,omitempty" gorm:"type:double precision"`
	High                 float64   `json:"high,omitempty" gorm:"type:double precision"`
	Low                  float64   `json:"low,omitempty" gorm:"type:double precision"`
	Close                float64   `json:"close,omitempty" gorm:"type:double precision"`
	Volume               int       `json:"volume,omitempty"`
	Average              float64   `json:"average,omitempty" gorm:"type:double precision"`
	Notional             float64   `json:"notional,omitempty" gorm:"type:double precision"`
	NumberOfTrades       int       `json:"numberOfTrades,omitempty"`
	ChangeOverTime       float64   `json:"changeOverTime,omitempty" gorm:"type:double precision"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date field, which is specified as "YYYY-MM-DD",
// and the Minute field, which is `/\d{2}:\d{2}/` into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled or if the date parsing fails.
func (i *Intraday) UnmarshalJSON(data []byte) (err error) {
	type intraday Intraday
	type embedded struct {
		intraday
		Date string `json:"date"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err != nil {
		return
	}
	*i = Intraday(tmp.intraday)
	var easternTime *time.Location
	if easternTime, err = time.LoadLocation("America/New_York"); err == nil {
		i.Date, err = time.ParseInLocation("2006-01-02T15:04", fmt.Sprintf("%sT%s", tmp.Date, tmp.Minute), easternTime)
	}
	return
}

// MarshalJSON satisfies the json.Marshaler interface.
// It undoes what UnmarshalJSON does.
func (i *Intraday) MarshalJSON() ([]byte, error) {
	type intraday Intraday
	type embedded struct {
		intraday
		Date string `json:"date"`
	}
	tmp := new(embedded)
	tmp.intraday = intraday(*i)
	tmp.Date = i.Date.Format("2006-01-02")
	return json.Marshal(tmp)
}

// Validate satisfies the Validator interface.
// It will return an error if the Date, Close or MarketClose fields are equal to their zero value.
func (i *Intraday) Validate() error {
	switch {
	case i.Date.IsZero():
		return fmt.Errorf("missing date")
	case i.MarketClose == 0:
		return fmt.Errorf("market close is zero")
	}
	return nil
}
