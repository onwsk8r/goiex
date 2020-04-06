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

package reference_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/reference"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Symbol", func() {
	var expected []Symbol

	BeforeEach(func() {
		expected = GoldenSymbol()
	})

	It("should parse symbols correctly", func() {
		var symbols []Symbol
		helper.TestdataFromJSON("core/reference/symbols.json", &symbols)
		Expect(symbols).To(BeEquivalentTo(expected))
	})

	Describe("Validate()", func() {
		var s Symbol
		BeforeEach(func() {
			s = expected[0]
		})

		It("should succeed if the Symbol is valid", func() {
			Expect(s.Validate()).To(Succeed())
		})
		It("should return an error if the Symbol is empty", func() {
			s.Symbol = ""
			Expect(s.Validate()).To(MatchError("missing symbol"))
		})
		It("should return an error if the IEXID is empty", func() {
			s.IEXID = ""
			Expect(s.Validate()).To(MatchError("missing IEX ID"))
		})
		It("should return an error if the Date is zero valued", func() {
			s.Date = time.Time{}
			Expect(s.Validate()).To(MatchError("missing date"))
		})
	})
})

var _ = XDescribe("Symbol Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []Symbol{
			Symbol{
				Symbol:   "A",
				Name:     "Agilent Technologies Inc.",
				Date:     time.Date(2019, time.March, 7, 0, 0, 0, 0, loc),
				Type:     "cs",
				IEXID:    "IEX_46574843354B2D52",
				Region:   "US",
				Currency: "USD",
				Enabled:  true,
			},
			Symbol{
				Symbol:   "AA",
				Name:     "Alcoa Corp.",
				Date:     time.Date(2019, time.March, 7, 0, 0, 0, 0, loc),
				Type:     "cs",
				IEXID:    "IEX_4238333734532D52",
				Region:   "US",
				Currency: "USD",
				Enabled:  true,
			},
		}
		helper.ToGolden("symbol", golden)
	})
})
