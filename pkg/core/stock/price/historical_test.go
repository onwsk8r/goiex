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

var _ = Describe("Historical", func() {
	var expected Historical
	BeforeEach(func() {
		expected = GoldenHistorical()
	})

	It("should parse historical prices correctly", func() {
		var res []Historical
		helper.TestdataFromJSON("core/stock/price/historical.json", &res)
		Expect(res[0]).To(Equal(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Historical is valid", func() {
			Expect(expected.Validate()).To(Succeed())
		})
		It("should return an error if the Date is zero valued", func() {
			expected.Date = time.Time{}
			Expect(expected.Validate()).To(MatchError("missing date"))
		})
		It("should return an error if the Close is zero", func() {
			expected.Close = 0
			Expect(expected.Validate()).To(MatchError("close is zero"))
		})
		It("should return an error if the UClose is zero", func() {
			expected.UClose = 0
			Expect(expected.Validate()).To(MatchError("unadjusted close is zero"))
		})
	})
})

var _ = XDescribe("Historical Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := Historical{
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
		}
		helper.ToGolden("historical", golden)
	})
})
