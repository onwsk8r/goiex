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
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// UpcomingDividend represents a data point from the upcoming events endpoint.
// https://iexcloud.io/docs/api/#upcoming-events
type UpcomingEarning struct {
	Symbol     string    `json:"symbol"`
	ReportDate time.Time `json:"reportDate"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the ReportDate field, which is specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (u *UpcomingEarning) UnmarshalJSON(data []byte) (err error) {
	type earning UpcomingEarning
	type embedded struct {
		earning
		ReportDate string `json:"reportDate"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		u.Symbol = tmp.Symbol
		u.ReportDate, _ = time.Parse("2006-01-02", tmp.ReportDate) // nolint: errcheck
		log.Debug().
			Dict("report_date", zerolog.Dict().Str("original", tmp.ReportDate).Time("parsed", u.ReportDate)).
			Msg("upcoming_earning: parsed report date")
	}
	return
}

// Validate satisfies the Validator interface.
// It will return an error if the ReportDate is zero or the Symbol is missing
func (u *UpcomingEarning) Validate() error {
	switch {
	case u.Symbol == "":
		return fmt.Errorf("symbol is missing")
	case u.ReportDate.IsZero():
		return fmt.Errorf("report date is missing")
	}
	return nil
}
