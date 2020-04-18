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

package market_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/market"
	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("UpcomingDividend", func() {
	var expected []UpcomingDividend
	BeforeEach(func() {
		expected = GoldenUpcomingDividends()
	})

	It("should parse upcoming dividends correctly", func() {
		var res []UpcomingDividend
		helper.TestdataFromJSON("core/market/upcoming_dividends.json", &res)
		Expect(res).To(ConsistOf(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the UpcomingDividend is valid", func() {
			for idx := range expected {
				Expect(expected[idx].Validate()).To(Succeed())
			}
		})
		It("should return an error if the Symbol is empty", func() {
			expected[0].Symbol = ""
			Expect(expected[0].Validate()).To(MatchError("symbol is missing"))
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected[0].ExDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("ex date is missing"))
		})
		It("should return an error if the DeclaredDate is zero valued", func() {
			expected[0].DeclaredDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("declared date is missing"))
		})
	})
})

var _ = XDescribe("UpcomingDividend Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []UpcomingDividend{
			UpcomingDividend{
				Dividend: fundamental.Dividend{
					Symbol:       "TUR",
					ExDate:       time.Date(2020, time.June, 16, 0, 0, 0, 0, loc),
					PaymentDate:  time.Date(2020, time.June, 26, 0, 0, 0, 0, loc),
					RecordDate:   time.Date(2020, time.June, 16, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2019, time.December, 20, 0, 0, 0, 0, loc),
					Amount:       0,
					Flag:         "",
					Currency:     "",
					Description:  "krT SEeSCFu CB Trse hiTyMIa",
					Frequency:    "",
				},
			},
			UpcomingDividend{
				Dividend: fundamental.Dividend{
					Symbol:       "GFY",
					ExDate:       time.Date(2020, time.February, 1, 0, 0, 0, 0, loc),
					PaymentDate:  time.Date(2020, time.February, 8, 0, 0, 0, 0, loc),
					RecordDate:   time.Date(2020, time.January, 26, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2019, time.November, 27, 0, 0, 0, 0, loc),
					Amount:       0.0799,
					Flag:         "aChs",
					Currency:     "USD",
					Description:  "anrSyrrO sdieha",
					Frequency:    "holnytm",
				},
			},
			UpcomingDividend{
				Dividend: fundamental.Dividend{
					Symbol:       "UACJF",
					ExDate:       time.Date(2020, time.April, 2, 0, 0, 0, 0, loc),
					RecordDate:   time.Date(2020, time.April, 15, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2019, time.May, 15, 0, 0, 0, 0, loc),
					Amount:       62,
					Flag:         "shaC",
					Currency:     "JPY",
					Description:  "eOrnrSsh adayri",
					Frequency:    "nea-unsmail",
				},
			},
		}
		helper.ToGolden("upcoming_dividends", golden)
	})
})
