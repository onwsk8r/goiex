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

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Dividend", func() {
	var expected []Dividend
	BeforeEach(func() {
		expected = []Dividend{{
			Amount:       0.70919585493507512,
			Currency:     "USD",
			DeclaredDate: time.Date(2020, time.October, 19, 0, 0, 0, 0, time.UTC),
			Description:  "Ordinary Shares",
			ExDate:       time.Date(2020, time.October, 28, 0, 0, 0, 0, time.UTC),
			Flag:         "Cash",
			Frequency:    "quarterly",
			PaymentDate:  time.Date(2020, time.November, 6, 0, 0, 0, 0, time.UTC),
			RecordDate:   time.Date(2020, time.October, 28, 0, 0, 0, 0, time.UTC),
			RefID:        2096218,
			Symbol:       "AAPL",
			ID:           "DIVIDENDS",
			Key:          "AAPL",
			Subkey:       "2053393",
			Date:         time.Date(2021, time.February, 3, 22, 42, 46, 191*1e6, time.UTC),
			Updated:      time.Date(2021, time.February, 3, 22, 42, 46, 191*1e6, time.UTC),
		}}
	})

	It("should parse dividends correctly", func() {
		var res []Dividend
		helper.TestdataFromJSON("core/stock/fundamental/dividends.json", &res)
		Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
	})

	It("should match the golden file", func() {
		golden := GoldenDividends()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("dividend", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})

	Describe("Validate()", func() {
		It("should succeed if the Dividend is valid", func() {
			Expect(expected[0].Validate()).To(Succeed())
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected[0].ExDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("ex date is missing"))
		})
		It("should return an error if the Amount is zero", func() {
			expected[0].Amount = 0
			Expect(expected[0].Validate()).To(MatchError("amount is missing"))
		})
		It("should return an error if the Currency is zero", func() {
			expected[0].Currency = ""
			Expect(expected[0].Validate()).To(MatchError("currency is missing"))
		})
	})
})
