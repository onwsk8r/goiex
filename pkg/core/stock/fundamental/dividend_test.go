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
	var expected []Dividend
	BeforeEach(func() {
		expected = GoldenDividends()
	})

	It("should parse dividends correctly", func() {
		var res []Dividend
		helper.TestdataFromJSON("core/stock/fundamental/dividends.json", &res)
		Expect(res).To(ConsistOf(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Dividend is valid", func() {
			for idx := range expected {
				Expect(expected[idx].Validate()).To(Succeed())
			}
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected[0].ExDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("ex date is missing"))
		})
		It("should return an error if the Amount is zero", func() {
			expected[1].Amount = 0
			Expect(expected[1].Validate()).To(MatchError("amount is missing"))
		})
		It("should return an error if the Currency is zero", func() {
			expected[2].Currency = ""
			Expect(expected[2].Validate()).To(MatchError("currency is missing"))
		})
	})
})

var _ = XDescribe("Dividend Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []Dividend{
			Dividend{
				ExDate:       time.Date(2020, time.January, 30, 0, 0, 0, 0, loc),
				PaymentDate:  time.Date(2020, time.February, 8, 0, 0, 0, 0, loc),
				RecordDate:   time.Date(2020, time.January, 24, 0, 0, 0, 0, loc),
				DeclaredDate: time.Date(2020, time.January, 10, 0, 0, 0, 0, loc),
				Amount:       0.41,
				Flag:         "Chas",
				Currency:     "USD",
				Description:  "Uitn",
				Frequency:    "uyarrleqt",
				Date:         time.Date(2020, time.April, 19, 0, 0, 0, 0, loc),
			},
			Dividend{
				ExDate:       time.Date(2019, time.October, 27, 0, 0, 0, 0, loc),
				PaymentDate:  time.Date(2019, time.November, 13, 0, 0, 0, 0, loc),
				RecordDate:   time.Date(2019, time.November, 2, 0, 0, 0, 0, loc),
				DeclaredDate: time.Date(2019, time.October, 16, 0, 0, 0, 0, loc),
				Amount:       0.39,
				Flag:         "Csha",
				Currency:     "USD",
				Description:  "Unti",
				Frequency:    "rqeyltaur",
				Date:         time.Date(2020, time.April, 19, 0, 0, 0, 0, loc),
			},
			Dividend{
				ExDate:       time.Date(2019, time.August, 6, 0, 0, 0, 0, loc),
				PaymentDate:  time.Date(2019, time.August, 11, 0, 0, 0, 0, loc),
				RecordDate:   time.Date(2019, time.July, 26, 0, 0, 0, 0, loc),
				DeclaredDate: time.Date(2019, time.July, 17, 0, 0, 0, 0, loc),
				Amount:       0.4,
				Flag:         "aCsh",
				Currency:     "USD",
				Description:  "iUtn",
				Frequency:    "raqletyru",
				Date:         time.Date(2020, time.April, 19, 0, 0, 0, 0, loc),
			},
			Dividend{
				ExDate:       time.Date(2019, time.April, 29, 0, 0, 0, 0, loc),
				PaymentDate:  time.Date(2019, time.April, 30, 0, 0, 0, 0, loc),
				RecordDate:   time.Date(2019, time.April, 23, 0, 0, 0, 0, loc),
				DeclaredDate: time.Date(2019, time.April, 13, 0, 0, 0, 0, loc),
				Amount:       0.41,
				Flag:         "hCas",
				Currency:     "USD",
				Description:  "Uint",
				Frequency:    "qrauetryl",
				Date:         time.Date(2020, time.April, 19, 0, 0, 0, 0, loc),
			},
		}
		helper.ToGolden("dividend", golden)
	})
})
