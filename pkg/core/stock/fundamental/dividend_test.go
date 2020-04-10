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

var _ = Describe("Dividend", func() {
	var expected Dividend
	BeforeEach(func() {
		expected = GoldenDividend()
	})

	It("should parse dividends correctly", func() {
		var res []Dividend
		helper.TestdataFromJSON("core/stock/fundamental/dividends.json", &res)
		Expect(res[0]).To(Equal(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Dividend is valid", func() {
			Expect(expected.Validate()).To(Succeed())
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected.ExDate = time.Time{}
			Expect(expected.Validate()).To(MatchError("ex date is missing"))
		})
		It("should return an error if the Amount is zero", func() {
			expected.Amount = 0
			Expect(expected.Validate()).To(MatchError("amount is missing"))
		})
		It("should return an error if the Currency is zero", func() {
			expected.Currency = ""
			Expect(expected.Validate()).To(MatchError("currency is missing"))
		})
	})
})

var _ = XDescribe("Dividend Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := Dividend{
			Symbol:       "AAPL",
			ExDate:       time.Date(2017, time.August, 10, 0, 0, 0, 0, loc),
			PaymentDate:  time.Date(2017, time.August, 17, 0, 0, 0, 0, loc),
			RecordDate:   time.Date(2017, time.August, 14, 0, 0, 0, 0, loc),
			DeclaredDate: time.Date(2017, time.August, 1, 0, 0, 0, 0, loc),
			Amount:       0.63,
			Flag:         "Dividend income",
			Currency:     "USD",
			Description:  "Apple declares dividend of .63",
			Frequency:    "quarterly",
		}
		helper.ToGolden("dividend", golden)
	})
})
