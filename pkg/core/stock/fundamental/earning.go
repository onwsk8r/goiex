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
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Earning represents a data point from the Earnings endpoint.
// https://iexcloud.io/docs/api/#earnings
type Earning struct {
	ActualEPS            float64   `json:"actualEPS"`
	ConsensusEPS         float64   `json:"consensusEPS"`
	AnnounceTime         string    `json:"announceTime"` // TODO: BTO, DMT, AMC
	NumberOfEstimates    int       `json:"numberOfEstimates"`
	EPSSurpriseDollar    float64   `json:"EPSSurpriseDollar"`
	EPSReportDate        time.Time `json:"EPSReportDate"`
	FiscalPeriod         string    `json:"fiscalPeriod"`
	FiscalEndDate        time.Time `json:"fiscalEndDate"`
	YearAgo              float64   `json:"yearAgo"`
	YearAgoChangePercent float64   `json:"yearAgoChangePercent"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// This function correctly translates the date fields, which are specified as "YYYY-MM-DD",
// into time.Times by using time.Parse().
// It will return an error if the JSON cannot be unmarshaled, but NOT if the date parsing fails.
func (e *Earning) UnmarshalJSON(data []byte) (err error) {
	type earning Earning
	type embedded struct {
		earning
		EPSReportDate string `json:"EPSReportDate"`
		FiscalEndDate string `json:"fiscalEndDate"`
	}
	tmp := new(embedded)
	if err = json.Unmarshal(data, tmp); err == nil {
		*e = Earning(tmp.earning)
		e.EPSReportDate, _ = time.Parse("2006-01-02", tmp.EPSReportDate) // nolint: errcheck
		e.FiscalEndDate, _ = time.Parse("2006-01-02", tmp.FiscalEndDate) // nolint: errcheck
		log.Debug().
			Dict("report_date", zerolog.Dict().Str("original", tmp.EPSReportDate).Time("parsed", e.EPSReportDate)).
			Dict("fiscal_end_date", zerolog.Dict().Str("original", tmp.FiscalEndDate).Time("parsed", e.FiscalEndDate)).
			Msg("earning: parsed dates")
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
