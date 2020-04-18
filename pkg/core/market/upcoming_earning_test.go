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
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("UpcomingEarning", func() { // nolint: dupl
	var expected []UpcomingEarning
	BeforeEach(func() {
		expected = GoldenUpcomingEarnings()
	})

	It("should parse upcoming earnings correctly", func() {
		var res []UpcomingEarning
		helper.TestdataFromJSON("core/market/upcoming_earnings.json", &res)
		Expect(res).To(ConsistOf(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the UpcomingEarning is valid", func() {
			for idx := range expected {
				Expect(expected[idx].Validate()).To(Succeed())
			}
		})
		It("should return an error if the Symbol is empty", func() {
			expected[0].Symbol = ""
			Expect(expected[0].Validate()).To(MatchError("symbol is missing"))
		})
		It("should return an error if the ReportDate is zero valued", func() {
			expected[0].ReportDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("report date is missing"))
		})
	})
})

var _ = XDescribe("UpcomingEarning Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []UpcomingEarning{
			UpcomingEarning{
				Symbol:     "RESN",
				ReportDate: time.Date(2020, time.March, 8, 0, 0, 0, 0, loc),
			},
			UpcomingEarning{
				Symbol:     "KHOLY",
				ReportDate: time.Date(2020, time.February, 21, 0, 0, 0, 0, loc),
			},
			UpcomingEarning{
				Symbol:     "SWM",
				ReportDate: time.Date(2020, time.March, 1, 0, 0, 0, 0, loc),
			},
		}
		helper.ToGolden("upcoming_earnings", golden)
	})
})
