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

package fundamental_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Earning", func() {
	var expected []Earning
	BeforeEach(func() {
		expected = GoldenEarnings()
	})

	It("should parse earnings correctly", func() {
		var res = struct {
			Symbol   string    `json:"symbol"`
			Earnings []Earning `json:"earnings"`
		}{}
		helper.TestdataFromJSON("core/stock/fundamental/earnings.json", &res)
		Expect(res.Earnings).To(ConsistOf(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Earning is valid", func() {
			Expect(expected[0].Validate()).To(Succeed())
		})
		It("should return an error if the ActualEPS is zero valued", func() {
			expected[0].ActualEPS = 0
			Expect(expected[0].Validate()).To(MatchError("actual EPS is zero"))
		})
		It("should return an error if the ReportDate is zero valued", func() {
			expected[0].EPSReportDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("report date is missing"))
		})
	})
})

var _ = Describe("Earning Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []Earning{
			Earning{
				ActualEPS:            2.46,
				ConsensusEPS:         2.36,
				AnnounceTime:         "AMC",
				NumberOfEstimates:    34,
				EPSSurpriseDollar:    0.1,
				EPSReportDate:        time.Date(2019, time.April, 30, 0, 0, 0, 0, loc),
				FiscalPeriod:         "Q1 2019",
				FiscalEndDate:        time.Date(2019, time.March, 31, 0, 0, 0, 0, loc),
				YearAgo:              2.73,
				YearAgoChangePercent: -0.0989,
			},
			Earning{
				ActualEPS:            4.18,
				ConsensusEPS:         4.17,
				AnnounceTime:         "AMC",
				NumberOfEstimates:    35,
				EPSSurpriseDollar:    0.01,
				EPSReportDate:        time.Date(2019, time.January, 29, 0, 0, 0, 0, loc),
				FiscalPeriod:         "Q4 2018",
				FiscalEndDate:        time.Date(2018, time.December, 31, 0, 0, 0, 0, loc),
				YearAgo:              3.89,
				YearAgoChangePercent: 0.0746,
			},
		}
		helper.ToGolden("earning", golden)
	})
})
