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

// Earning represents a data point from the Earnings endpoint.
// https://iexcloud.io/docs/api/#earnings
type Earning struct {
	EPSReportDate            time.Time `json:"EPSReportDate,omitempty"`
	EPSSurpriseDollar        float64   `json:"EPSSurpriseDollar,omitempty"`
	EPSSurpriseDollarPercent float64   `json:"EPSSurpriseDollarPercent,omitempty"`
	ActualEPS                float64   `json:"actualEPS,omitempty"`
	AnnounceTime             string    `json:"announceTime,omitempty"` // TODO: BTO, DMT, AMC
	ConsensusEPS             float64   `json:"consensusEPS,omitempty"`
	Currency                 string    `json:"currency,omitempty"`
	FiscalEndDate            time.Time `json:"fiscalEndDate,omitempty"`
	FiscalPeriod             string    `json:"fiscalPeriod,omitempty"`
	NumberOfEstimates        int       `json:"numberOfEstimates,omitempty"`
	PeriodType               string    `json:"periodType,omitempty"`
	Symbol                   string    `json:"symbol,omitempty"`
	YearAgo                  float64   `json:"yearAgo,omitempty"`
	YearAgoChangePercent     float64   `json:"yearAgoChangePercent,omitempty"`
	ID                       string    `json:"id,omitempty"`
	Source                   string    `json:"source,omitempty"`
	Key                      string    `json:"key,omitempty"`
	Subkey                   string    `json:"subkey,omitempty"`
	Date                     time.Time `json:"date,omitempty"`
	Updated                  time.Time `json:"updated,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (e *Earning) UnmarshalJSON(data []byte) (err error) { // nolint:dupl
	type earning Earning
	type embedded struct {
		earning
		EPSReportDate string `json:"EPSReportDate"`
		FiscalEndDate string `json:"fiscalEndDate"`
		Date          int64  `json:"date,omitempty"`
		Updated       int64  `json:"updated,omitempty"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*e = Earning(tmp.earning)
		e.EPSReportDate, _ = time.Parse("2006-01-02", tmp.EPSReportDate) // nolint:errcheck
		e.FiscalEndDate, _ = time.Parse("2006-01-02", tmp.FiscalEndDate) // nolint:errcheck
		e.Date = time.Unix(tmp.Date/1000, tmp.Date%1000*1e6)             // nolint:gomnd
		e.Updated = time.Unix(tmp.Updated/1000, tmp.Updated%1000*1e6)    // nolint:gomnd
		log.Debug().Interface("original", tmp).Interface("final", e).Msg("earning: parsed dates")
	}
	return
}

// Validate satisfies the Validator interface.
// It will return an error if the ActualEPS, ConsensusEPS, or EPSReportDate fields are zero
func (e *Earning) Validate() error {
	switch {
	case e.ActualEPS == 0:
		return fmt.Errorf("actual EPS is zero")
	case e.EPSReportDate.IsZero():
		return fmt.Errorf("report date is missing")
	}
	return nil
}
