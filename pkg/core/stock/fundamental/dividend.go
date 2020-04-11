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
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Dividend represents a data point from the basic dividends endpoint.
// https://iexcloud.io/docs/api/#dividends-basic
type Dividend struct {
	Symbol       string    `json:"symbol"`
	ExDate       time.Time `json:"exDate"`
	PaymentDate  time.Time `json:"paymentDate"`
	RecordDate   time.Time `json:"recordDate"`
	DeclaredDate time.Time `json:"declaredDate"`
	Amount       float64   `json:"amount"`
	Flag         string    `json:"flag"`
	Currency     string    `json:"currency"`
	Description  string    `json:"description"`
	Frequency    string    `json:"frequency"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (d *Dividend) UnmarshalJSON(data []byte) (err error) {
	type dividend Dividend
	type embedded struct {
		dividend
		Amount       string `json:"amount"`
		ExDate       string `json:"exDate"`
		PaymentDate  string `json:"paymentDate"`
		RecordDate   string `json:"recordDate"`
		DeclaredDate string `json:"declaredDate"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*d = Dividend(tmp.dividend)
		// Ignore date parsing issues in case one or more dates are missing
		d.ExDate, _ = time.Parse("2006-01-02", tmp.ExDate)             // nolint: errcheck
		d.PaymentDate, _ = time.Parse("2006-01-02", tmp.PaymentDate)   // nolint: errcheck
		d.RecordDate, _ = time.Parse("2006-01-02", tmp.RecordDate)     // nolint: errcheck
		d.DeclaredDate, _ = time.Parse("2006-01-02", tmp.DeclaredDate) // nolint: errcheck
		d.Amount, _ = strconv.ParseFloat(tmp.Amount, 64)               // nolint: errcheck
		log.Debug().
			Dict("ex_date", zerolog.Dict().Str("original", tmp.ExDate).Time("parsed", d.ExDate)).
			Dict("payment_date", zerolog.Dict().Str("original", tmp.PaymentDate).Time("parsed", d.PaymentDate)).
			Dict("record_date", zerolog.Dict().Str("original", tmp.RecordDate).Time("parsed", d.RecordDate)).
			Dict("declared_date", zerolog.Dict().Str("original", tmp.DeclaredDate).Time("parsed", d.DeclaredDate)).
			Dict("amount", zerolog.Dict().Str("original", tmp.Amount).Float64("parsed", d.Amount)).
			Msg("dividend: parsed date")
	}
	return
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
