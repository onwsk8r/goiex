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

var _ = Describe("Intraday", func() { // nolint: dupl
	var expected []Intraday
	BeforeEach(func() {
		expected = []Intraday{{
			Date:                 time.Date(2017, 12, 15, 14, 30, 0, 0, time.UTC),
			Minute:               "09:30",
			Label:                "09:30 AM",
			MarketOpen:           143.98,
			MarketClose:          143.775,
			MarketHigh:           143.98,
			MarketLow:            143.775,
			MarketAverage:        143.889,
			MarketVolume:         3070,
			MarketNotional:       441740.275,
			MarketNumberOfTrades: 20,
			MarketChangeOverTime: -0.004,
			High:                 143.98,
			Low:                  143.775,
			Open:                 143.98,
			Close:                143.775,
			Average:              143.889,
			Volume:               3070,
			Notional:             441740.275,
			NumberOfTrades:       20,
			ChangeOverTime:       -0.0039,
		}}
	})

	It("should parse intraday prices correctly", func() {
		var res []Intraday
		helper.TestdataFromJSON("core/stock/intraday.json", &res)
		Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
	})

	It("should match the golden file", func() {
		golden := GoldenIntraday()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("intraday", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})

	Describe("Validate()", func() {
		It("should succeed if the Intraday is valid", func() {
			Expect(expected[0].Validate()).To(Succeed())
		})
		It("should return an error if the Date is zero valued", func() {
			expected[0].Date = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("missing date"))
		})
		It("should return an error if the MarketClose is zero", func() {
			expected[0].MarketClose = 0
			Expect(expected[0].Validate()).To(MatchError("market close is zero"))
		})
	})
})
