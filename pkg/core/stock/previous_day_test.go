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

var _ = Describe("PreviousDay", func() {
	var expected PreviousDay
	BeforeEach(func() {
		expected = PreviousDay{Historical{
			Close:   116.59,
			High:    117.49,
			Low:     116.22,
			Open:    116.57,
			Symbol:  "AAPL",
			Volume:  46691331,
			ID:      "HISTORICAL_PRICES",
			Key:     "AAPL",
			Date:    time.Date(2020, time.November, 30, 14, 33, 10, 0, time.UTC),
			Updated: time.Date(2020, time.November, 30, 14, 33, 10, 0, time.UTC),
			UOpen:   116.57,
			UHigh:   117.49,
			ULow:    116.22,
			UClose:  116.59,
			UVolume: 46691331,
			FOpen:   116.57,
			FHigh:   117.49,
			FLow:    116.22,
			FClose:  116.59,
			FVolume: 46691331,
			Label:   "Nov 27, 20",
		}}
	})

	It("should parse previous day prices correctly", func() {
		var res PreviousDay
		helper.TestdataFromJSON("core/stock/previous_day.json", &res)
		Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
	})

	It("should match the golden file", func() {
		golden := GoldenPreviousDay()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("previous_day", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})

	Describe("Validate()", func() {
		It("should succeed if the PreviousDay is valid", func() {
			Expect(expected.Validate()).To(Succeed())
		})
		It("should return an error if the Symbol is empty", func() {
			expected.Symbol = ""
			Expect(expected.Validate()).To(MatchError("missing symbol"))
		})
		It("should return an error if the Historical is invalid", func() {
			expected.Date = time.Time{}
			Expect(expected.Validate()).ToNot(Succeed())
		})
	})
})
