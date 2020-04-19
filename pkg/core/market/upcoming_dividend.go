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

package market

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// UpcomingDividend represents a data point from the upcoming events endpoint.
// The data structure is similar to the basic dividends endpoint with an added symbol field. See
// https://iexcloud.io/docs/api/#dividends-basic and
// https://iexcloud.io/docs/api/#upcoming-events for more information.
type UpcomingDividend struct {
	ExDate       time.Time `json:"exDate"`
	PaymentDate  time.Time `json:"paymentDate"`
	RecordDate   time.Time `json:"recordDate"`
	DeclaredDate time.Time `json:"declaredDate"`
	Amount       float64   `json:"amount"`
	Flag         string    `json:"flag"`
	Currency     string    `json:"currency"`
	Description  string    `json:"description"`
	Frequency    string    `json:"frequency"`
	Symbol       string    `json:"symbol"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (u *UpcomingDividend) UnmarshalJSON(data []byte) (err error) {
	type dividend UpcomingDividend
	type embedded struct {
		dividend
		Amount       string `json:"amount"`
		ExDate       string `json:"exDate"`
		PaymentDate  string `json:"paymentDate"`
		RecordDate   string `json:"recordDate"`
		DeclaredDate string `json:"declaredDate"`
	}
	tmp := new(embedded)
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}
	*u = UpcomingDividend(tmp.dividend)
	// Ignore date parsing issues in case one or more dates are missing
	u.ExDate, _ = time.Parse("2006-01-02", tmp.ExDate)             // nolint: errcheck
	u.PaymentDate, _ = time.Parse("2006-01-02", tmp.PaymentDate)   // nolint: errcheck
	u.RecordDate, _ = time.Parse("2006-01-02", tmp.RecordDate)     // nolint: errcheck
	u.DeclaredDate, _ = time.Parse("2006-01-02", tmp.DeclaredDate) // nolint: errcheck
	u.Amount, _ = strconv.ParseFloat(tmp.Amount, 64)               // nolint: errcheck
	log.Debug().
		Dict("ex_date", zerolog.Dict().Str("original", tmp.ExDate).Time("parsed", u.ExDate)).
		Dict("payment_date", zerolog.Dict().Str("original", tmp.PaymentDate).Time("parsed", u.PaymentDate)).
		Dict("record_date", zerolog.Dict().Str("original", tmp.RecordDate).Time("parsed", u.RecordDate)).
		Dict("declared_date", zerolog.Dict().Str("original", tmp.DeclaredDate).Time("parsed", u.DeclaredDate)).
		Dict("amount", zerolog.Dict().Str("original", tmp.Amount).Float64("parsed", u.Amount)).
		Msg("upcoming_dividend: parsed date")
	return nil
}

// Validate satisfies the Validator interface.
// It will return an error if the Symbol or ex-date are zero-valued. It will not
// validate the emdedded dividend because required values - such as the amount - may be missing.
func (u *UpcomingDividend) Validate() (err error) {
	switch {
	case u.Symbol == "":
		err = fmt.Errorf("symbol is missing")
	case u.ExDate.IsZero():
		err = fmt.Errorf("ex date is missing")
	}
	return err
}
