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

package stock_test

import (
	"time"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Earning", func() {
	var expected []Earning
	BeforeEach(func() {
		expected = []Earning{{
			EPSReportDate:            time.Date(2020, time.October, 22, 0, 0, 0, 0, time.UTC),
			EPSSurpriseDollar:        0.330098697767743855,
			EPSSurpriseDollarPercent: 0.1440300992341432,
			ActualEPS:                func(i float64) *float64 { return &i }(5.07),
			AnnounceTime:             "AMC",
			ConsensusEPS:             func(i float64) *float64 { return &i }(2.703327),
			Currency:                 "USD",
			FiscalEndDate:            time.Date(2020, time.September, 18, 0, 0, 0, 0, time.UTC),
			FiscalPeriod:             "Q4 2020",
			NumberOfEstimates:        31,
			PeriodType:               "quarterly",
			Symbol:                   "AAPL",
			YearAgo:                  func(i float64) *float64 { return &i }(4.7006),
			YearAgoChangePercent:     -0.7303425659240847,
			ID:                       "PREMIUM_EARNINGS",
			Key:                      "AAPL",
			Subkey:                   "Q42020",
			Date:                     time.Date(2020, time.November, 25, 13, 48, 45, 71*1e6, time.UTC),
			Updated:                  time.Date(2020, time.November, 25, 13, 48, 45, 71*1e6, time.UTC),
		}}
	})

	It("should parse earnings correctly", func() {
		var res = struct {
			Symbol   string    `json:"symbol"`
			Earnings []Earning `json:"earnings"`
		}{}
		helper.TestdataFromJSON("core/stock/earnings.json", &res)
		Expect(cmp.Equal(expected, res.Earnings)).To(BeTrue(), cmp.Diff(expected, res.Earnings))
	})

	It("should match the golden file", func() {
		golden := GoldenEarnings()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("earning", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})

	Describe("Validate()", func() {
		It("should succeed if the Earning is valid", func() {
			Expect(expected[0].Validate()).To(Succeed())
		})
		It("should return an error if the ActualEPS is zero valued", func() {
			expected[0].ActualEPS = nil
			Expect(expected[0].Validate()).To(MatchError("actual EPS is zero"))
		})
		It("should return an error if the ReportDate is zero valued", func() {
			expected[0].EPSReportDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("report date is missing"))
		})
	})
})
