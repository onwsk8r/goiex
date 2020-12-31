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

// Dividend represents a data point from the basic dividends endpoint.
// Note the refid field exists in the sample data but not in the schema docs.
// https://iexcloud.io/docs/api/#dividends-basic
type Dividend struct {
	Symbol       string    `json:"symbol,omitempty" gorm:"primaryKey;type:character varying"`
	Amount       float64   `json:"amount,omitempty" gorm:"type:double precision"`
	Currency     string    `json:"currency,omitempty" gorm:"type:character varying"`
	DeclaredDate time.Time `json:"declaredDate,omitempty" gorm:"type:date"`
	Description  string    `json:"description,omitempty" gorm:"type:character varying"`
	ExDate       time.Time `json:"exDate,omitempty" gorm:"primaryKey;type:date"`
	Flag         string    `json:"flag,omitempty" gorm:"type:character varying"`
	Frequency    string    `json:"frequency,omitempty" gorm:"type:character varying"`
	PaymentDate  time.Time `json:"paymentDate,omitempty" gorm:"type:date"`
	RecordDate   time.Time `json:"recordDate,omitempty" gorm:"type:date"`
	RefID        float64   `json:"refid,omitempty" gorm:"-"`
	ID           string    `json:"id,omitempty" gorm:"-"`
	Source       string    `json:"source,omitempty" gorm:"-"`
	Key          string    `json:"key,omitempty" gorm:"-"`
	Subkey       string    `json:"subkey,omitempty" gorm:"-"`
	Date         time.Time `json:"date,omitempty" gorm:"type:date"`
	Updated      time.Time `json:"updated,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (d *Dividend) UnmarshalJSON(data []byte) (err error) {
	type dividend Dividend
	type embedded struct {
		dividend
		DeclaredDate string `json:"declaredDate,omitempty"`
		ExDate       string `json:"exDate,omitempty"`
		PaymentDate  string `json:"paymentDate,omitempty"`
		RecordDate   string `json:"recordDate,omitempty"`
		Date         int64  `json:"date,omitempty"`
		Updated      int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}
	*d = Dividend(tmp.dividend)
	// Ignore date parsing issues in case one or more dates are missing
	d.DeclaredDate, _ = time.Parse("2006-01-02", tmp.DeclaredDate) // nolint:errcheck
	d.ExDate, _ = time.Parse("2006-01-02", tmp.ExDate)             // nolint:errcheck
	d.PaymentDate, _ = time.Parse("2006-01-02", tmp.PaymentDate)   // nolint:errcheck
	d.RecordDate, _ = time.Parse("2006-01-02", tmp.RecordDate)     // nolint:errcheck
	d.Date = time.Unix(tmp.Date/1000, tmp.Date%1000*1e6)           // nolint:gomnd
	d.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6)  // nolint:gomnd
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
// It undoes what UnmarshalJSON does.
func (d *Dividend) MarshalJSON() ([]byte, error) {
	type dividend Dividend
	type embedded struct {
		dividend
		DeclaredDate string `json:"declaredDate,omitempty"`
		ExDate       string `json:"exDate,omitempty"`
		PaymentDate  string `json:"paymentDate,omitempty"`
		RecordDate   string `json:"recordDate,omitempty"`
		Date         int64  `json:"date,omitempty"`
		Updated      int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	tmp.dividend = dividend(*d)
	tmp.DeclaredDate = d.DeclaredDate.Format("2006-01-02")
	tmp.ExDate = d.ExDate.Format("2006-01-02")
	tmp.PaymentDate = d.PaymentDate.Format("2006-01-02")
	tmp.RecordDate = d.RecordDate.Format("2006-01-02")
	tmp.Date = d.Date.UnixNano() / 1e6       // nolint:gomnd
	tmp.Updated = d.Updated.UnixNano() / 1e6 // nolint:gomnd
	return json.Marshal(tmp)
}

// Validate satisfies the Validator interface.
// It will return an error if the Amount, ExDate, or Currency fields are equal to their zero value.
func (d *Dividend) Validate() error {
	switch {
	case d.Amount == 0:
		return fmt.Errorf("amount is missing")
	case d.ExDate.IsZero():
		return fmt.Errorf("ex date is missing")
	case d.Currency == "":
		return fmt.Errorf("currency is missing")
	}
	return nil
}
