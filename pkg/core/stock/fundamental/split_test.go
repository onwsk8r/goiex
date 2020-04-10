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

var _ = Describe("Split", func() {
	var expected Split
	BeforeEach(func() {
		expected = GoldenSplit()
	})

	It("should parse splits correctly", func() {
		var res []Split
		helper.TestdataFromJSON("core/stock/fundamental/splits.json", &res)
		Expect(res[0]).To(Equal(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Split is valid", func() {
			Expect(expected.Validate()).To(Succeed())
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected.ExDate = time.Time{}
			Expect(expected.Validate()).To(MatchError("ex date is missing"))
		})
		It("should return an error if the ToFactor is not positive", func() {
			expected.ToFactor = 0
			Expect(expected.Validate()).To(MatchError("to factor is not positive"))
		})
		It("should return an error if the FromFactor is not positive", func() {
			expected.FromFactor = -4
			Expect(expected.Validate()).To(MatchError("from factor is not positive"))
		})
	})
})

var _ = XDescribe("Split Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := Split{
			ExDate:       time.Date(2017, time.August, 10, 0, 0, 0, 0, loc),
			DeclaredDate: time.Date(2017, time.August, 1, 0, 0, 0, 0, loc),
			Ratio:        0.142857,
			ToFactor:     7,
			FromFactor:   1,
			Description:  "7-for-1 split",
		}
		helper.ToGolden("split", golden)
	})
})
