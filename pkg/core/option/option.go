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

	"github.com/rs/zerolog/log"
)

// Option represents a data point from the Options EOD endpoint.
// https://iexcloud.io/docs/api/#end-of-day-options
// Of all the date types, LastUpdated looks like the last day (ie date)
// the price was updated, where Updated looks to be the timestamp of the
// last release, which may be the next morning after LastUpdated
type Option struct {
	Ask                 *float64  `json:"ask,omitempty"`
	Bid                 *float64  `json:"bid,omitempty"`
	CFICode             string    `json:"cfiCode,omitempty"`
	Close               *float64  `json:"close,omitempty"`
	ClosingPrice        *float64  `json:"closingPrice,omitempty"`
	ContractDescription string    `json:"contractDescription,omitempty"`
	ContractName        string    `json:"contractName,omitempty"`
	ContractSize        float64   `json:"contractSize,omitempty"`
	Currency            string    `json:"currency,omitempty"`
	ExchangeCode        string    `json:"exchangeCode,omitempty"`
	ExchangeMIC         string    `json:"exchangeMIC,omitempty"`
	ExerciseStyle       string    `json:"exerciseStyle,omitempty"`
	ExpirationDate      time.Time `json:"-"`
	FIGI                string    `json:"figi,omitempty"`
	High                *float64  `json:"high,omitempty"`
	IsAdjusted          bool      `json:"isAdjusted,omitempty"`
	LastTrade           time.Time `json:"-"`
	LastUpdated         time.Time `json:"lastUpdated,omitempty"`
	Low                 *float64  `json:"low,omitempty"`
	MarginPrice         *float64  `json:"marginPrice,omitempty"`
	Open                *float64  `json:"open,omitempty"`
	OpenInterest        *uint     `json:"openInterest,omitempty"`
	SettlementPrice     *float64  `json:"settlementPrice,omitempty"`
	Side                string    `json:"side,omitempty"`
	StrikePrice         float64   `json:"strikePrice,omitempty"`
	Symbol              string    `json:"symbol,omitempty"`
	Type                string    `json:"type,omitempty"`
	Volume              *uint     `json:"volume,omitempty"`
	ID                  string    `json:"id,omitempty"`
	Key                 string    `json:"key,omitempty"`
	Subkey              string    `json:"subkey,omitempty"`
	Date                time.Time `json:"-"`
	Updated             time.Time `json:"-"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the expirationDate and lastUpdated fields, which are specified as
// "YYYYMMDD" and "YYYY-MM-DD", respectively, into a time.Time by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the lastUpdated parsing fails.
func (o *Option) UnmarshalJSON(data []byte) (err error) {
	type option Option
	type embedded struct {
		option
		ExpirationDate string `json:"expirationDate,omitempty"`
		LastTradeDate  string `json:"lastTradeDate,omitempty"`
		LastTradeTime  string `json:"lastTradeTime,omitempty"`
		LastUpdated    string `json:"lastUpdated,omitempty"`
		Date           int64  `json:"date,omitempty"`
		Updated        int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err != nil {
		return
	}
	*o = Option(tmp.option)
	if o.ExpirationDate, err = time.Parse("20060102", tmp.ExpirationDate); err != nil {
		return
	}
	if o.LastUpdated, err = time.Parse("2006-01-02", tmp.LastUpdated); err != nil {
		return
	}
	val := fmt.Sprintf("%sT%s", tmp.LastTradeDate, tmp.LastTradeTime)
	o.LastTrade, _ = time.ParseInLocation("2006-01-02T15:04:05", val, time.UTC) // nolint:errcheck
	o.Date = time.Unix(tmp.Date/1000, tmp.Date%1000*1e6)                        // nolint:gomnd
	o.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6)               // nolint:gomnd
	log.Debug().Interface("original", tmp).Interface("final", o).Msg("option: parsed dates")
	return
}

func (o *Option) MarshalJSON() ([]byte, error) {
	type option Option
	type embedded struct {
		option
		ExpirationDate string `json:"expirationDate,omitempty"`
		LastTradeDate  string `json:"lastTradeDate,omitempty"`
		LastTradeTime  string `json:"lastTradeTime,omitempty"`
		LastUpdated    string `json:"lastUpdated,omitempty"`
		Date           int64  `json:"date,omitempty"`
		Updated        int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	tmp.option = option(*o)
	tmp.ExpirationDate = o.ExpirationDate.Format("20060102")
	tmp.LastUpdated = o.LastUpdated.Format("2006-01-02")
	tmp.LastTradeDate = o.LastTrade.Format("2006-01-02")
	tmp.LastTradeTime = o.LastTrade.Format("15:04:05")
	tmp.Date = o.Date.UnixNano() / 1e6       // nolint: gomnd
	tmp.Updated = o.Updated.UnixNano() / 1e6 // nolint: gomnd
	return json.Marshal(tmp)
}

// Validate satisfies the Validator interface.
// It will return an error if the Date, Close or UClose fields are equal to their zero value.
func (o *Option) Validate() error {
	switch {
	case o.LastUpdated.IsZero():
		return fmt.Errorf("missing last updated")
	case o.Symbol == "":
		return fmt.Errorf("missing symbol")
	case o.ID == "":
		return fmt.Errorf("missing id")
	case o.ExpirationDate.IsZero():
		return fmt.Errorf("missing expiration date")
	case o.StrikePrice == 0:
		return fmt.Errorf("strike price is zero")
	case o.Side == "":
		return fmt.Errorf("missing option side")
	}
	return nil
}
