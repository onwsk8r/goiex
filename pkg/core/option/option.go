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
)

// Option represents a data point from the Options EOD endpoint.
// https://iexcloud.io/docs/api/#end-of-day-options
// Of all the date types, LastUpdated looks like the last day (ie date)
// the price was updated, where Updated looks to be the timestamp of the
// last release, which may be the next morning after LastUpdated
type Option struct {
	Ask                 *float64  `json:"ask,omitempty" gorm:"type:double precision"`
	Bid                 *float64  `json:"bid,omitempty" gorm:"type:double precision"`
	CFICode             string    `json:"cfiCode,omitempty" gorm:"type:character varying"`
	Close               *float64  `json:"close,omitempty" gorm:"type:double precision"`
	ClosingPrice        *float64  `json:"closingPrice,omitempty" gorm:"type:double precision"`
	ContractDescription string    `json:"contractDescription,omitempty" gorm:"type:character varying"`
	ContractName        string    `json:"contractName,omitempty" gorm:"type:character varying"`
	ContractSize        float64   `json:"contractSize,omitempty" gorm:"type:double precision"`
	Currency            string    `json:"currency,omitempty" gorm:"type:character varying"`
	ExchangeCode        string    `json:"exchangeCode,omitempty" gorm:"type:character varying"`
	ExchangeMIC         string    `json:"exchangeMIC,omitempty" gorm:"type:character varying"`
	ExerciseStyle       string    `json:"exerciseStyle,omitempty" gorm:"type:character varying"`
	FIGI                string    `json:"figi,omitempty" gorm:"type:character varying"`
	High                *float64  `json:"high,omitempty" gorm:"type:double precision"`
	LastTrade           time.Time `json:"-"`
	LastUpdated         time.Time `json:"lastUpdated,omitempty" gorm:"primaryKey;type:date"`
	Low                 *float64  `json:"low,omitempty" gorm:"type:double precision"`
	MarginPrice         *float64  `json:"marginPrice,omitempty" gorm:"type:double precision"`
	Open                *float64  `json:"open,omitempty" gorm:"type:double precision"`
	OpenInterest        *uint     `json:"openInterest,omitempty" gorm:"type:double precision"`
	SettlementPrice     *float64  `json:"settlementPrice,omitempty" gorm:"type:double precision"`
	Symbol              string    `json:"symbol,omitempty" gorm:"primaryKey;type:character varying"`
	ExpirationDate      time.Time `json:"-" gorm:"primaryKey;type:date"`
	StrikePrice         float64   `json:"strikePrice,omitempty" gorm:"primaryKey;type:double precision"`
	Side                string    `json:"side,omitempty" gorm:"primaryKey;type:character varying"`
	IsAdjusted          bool      `json:"isAdjusted,omitempty" gorm:"primaryKey"`
	Type                string    `json:"type,omitempty" gorm:"type:character varying"`
	Volume              *uint     `json:"volume,omitempty"`
	ID                  string    `json:"id,omitempty" gorm:"-"`
	Key                 string    `json:"key,omitempty" gorm:"-"`
	Subkey              string    `json:"subkey,omitempty" gorm:"-"`
	Date                time.Time `json:"-" gorm:"type:date"`
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
