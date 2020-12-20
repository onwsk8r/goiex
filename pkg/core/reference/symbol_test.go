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

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/reference"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Symbol", func() {
	var expected []Symbol

	BeforeEach(func() {
		expected = []Symbol{{
			Symbol:   "A",
			Name:     "Agilent Technologies Inc.",
			Date:     time.Date(2019, time.March, 7, 0, 0, 0, 0, time.UTC),
			Type:     "cs",
			IEXID:    "IEX_46574843354B2D52",
			Region:   "US",
			Currency: "USD",
			Enabled:  true,
			FIGI:     "BBG000C2V3D6",
			CIK:      "1090872",
		}, {
			Symbol:   "AA",
			Name:     "Alcoa Corp.",
			Date:     time.Date(2019, time.March, 7, 0, 0, 0, 0, time.UTC),
			Type:     "cs",
			IEXID:    "IEX_4238333734532D52",
			Region:   "US",
			Currency: "USD",
			Enabled:  true,
			FIGI:     "BBG00B3T3HD3",
			CIK:      "1675149",
		},
		}
	})

	It("should parse symbols correctly", func() {
		var symbols []Symbol
		helper.TestdataFromJSON("core/reference/symbols.json", &symbols)
		Expect(symbols).To(BeEquivalentTo(expected), cmp.Diff(symbols, expected))
	})

	It("should match the golden file", func() {
		golden := GoldenSymbol()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("symbol", expected)
			Fail(cmp.Diff(golden, expected))
		}
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
