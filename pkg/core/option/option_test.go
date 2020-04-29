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

package option_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/option"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Option", func() {
	var expected Option
	BeforeEach(func() {
		expected = GoldenOption()
	})

	It("should parse option prices correctly", func() {
		var res []Option
		helper.TestdataFromJSON("core/option/options.json", &res)

		// See https://github.com/golang/go/issues/10089#issuecomment-77463157 for why this is necessary
		// My guess is because the "values pointed to" for the location are identical, but each time.Time
		// contains a different pointer (ie living in a different memory address).
		// Expect(res[0].ExpirationDate.Equal(expected.ExpirationDate)).To(BeTrue(), "expiration dates are inequal")
		// res[0].ExpirationDate = expected.ExpirationDate
		// Expect(res[0].LastUpdated.Equal(expected.LastUpdated)).To(BeTrue(), "last updated dates are inequal")
		// res[0].LastUpdated = expected.LastUpdated
		Expect(res[0]).To(Equal(expected), "values are inequal")
	})

	Describe("Validate()", func() {
		It("should succeed if the Option is valid", func() {
			Expect(expected.Validate()).To(Succeed())
		})
		It("should return an error if the Symbol is empty", func() {
			expected.Symbol = ""
			Expect(expected.Validate()).To(MatchError("missing symbol"))
		})
		It("should return an error if the ID is empty", func() {
			expected.ID = ""
			Expect(expected.Validate()).To(MatchError("missing id"))
		})
		It("should return an error if the ExpirationDate is zero valued", func() {
			expected.ExpirationDate = time.Time{}
			Expect(expected.Validate()).To(MatchError("missing expiration date"))
		})
		It("should return an error if the StrikePrice is zero", func() {
			expected.StrikePrice = 0
			Expect(expected.Validate()).To(MatchError("strike price is zero"))
		})
	})
})

var _ = XDescribe("Option Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := Option{
			Symbol:         "AAPL",
			ID:             "AAPL20190621C00240000",
			ExpirationDate: time.Date(2019, time.June, 21, 0, 0, 0, 0, loc),
			ContractSize:   100,
			StrikePrice:    240,
			ClosingPrice:   0.39,
			Side:           "call",
			Type:           "equity",
			Volume:         884,
			OpenInterest:   12197,
			Bid:            0.38,
			Ask:            0.42,
			LastUpdated:    time.Date(2019, time.April, 25, 0, 0, 0, 0, loc),
			IsAdjusted:     false,
		}
		helper.ToGolden("option", golden)
	})
})
