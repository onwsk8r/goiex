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

package price_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock/price"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("PreviousDay", func() {
	var expected PreviousDay
	BeforeEach(func() {
		expected = GoldenPreviousDay()
	})

	It("should parse previous day prices correctly", func() {
		var res PreviousDay
		helper.TestdataFromJSON("core/stock/price/previous_day.json", &res)
		Expect(res).To(Equal(expected))
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

var _ = XDescribe("PreviousDay Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := PreviousDay{
			Symbol: "AAPL",
			Historical: Historical{
				Date:           time.Date(2017, time.April, 3, 0, 0, 0, 0, loc),
				Open:           143.1192,
				High:           143.5275,
				Low:            142.4619,
				Close:          143.1092,
				Volume:         19985714,
				UOpen:          143.1192,
				UHigh:          143.5275,
				ULow:           142.4619,
				UClose:         143.1092,
				UVolume:        19985714,
				Change:         0.039835,
				ChangePercent:  0.028,
				Label:          "Apr 03, 17",
				ChangeOverTime: -0.0039,
			},
		}
		helper.ToGolden("previous_day", golden)
	})
})
