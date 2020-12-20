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

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/option"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Option", func() {
	var expected []Option
	BeforeEach(func() {
		expected = []Option{{
			Ask:                 func(i float64) *float64 { return &i }(0.1),
			Bid:                 func(i float64) *float64 { return &i }(0.05),
			CFICode:             "OPAXXX",
			Close:               func(i float64) *float64 { return &i }(0.05),
			ClosingPrice:        func(i float64) *float64 { return &i }(0.05),
			ContractDescription: "CVM Option Put 18/12/2020 2.5 on Ordinary Shares",
			ContractName:        "Cel-Sci Corp",
			ContractSize:        100,
			Currency:            "USD",
			ExerciseStyle:       "A",
			ExpirationDate:      time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
			FIGI:                "BBG00PZWYFY7",
			High:                func(i float64) *float64 { return &i }(0.1),
			IsAdjusted:          false,
			LastTrade:           time.Date(2020, 11, 25, 20, 24, 13, 0, time.UTC),
			LastUpdated:         time.Date(2020, 11, 29, 0, 0, 0, 0, time.UTC),
			Low:                 func(i float64) *float64 { return &i }(0.05),
			MarginPrice:         func(i float64) *float64 { return &i }(-0),
			Open:                func(i float64) *float64 { return &i }(0.1),
			OpenInterest:        func(i uint) *uint { return &i }(2039),
			SettlementPrice:     func(i float64) *float64 { return &i }(-0),
			Side:                "put",
			StrikePrice:         3,
			Symbol:              "CVM",
			Type:                "equity",
			Volume:              func(i uint) *uint { return &i }(17),
			ID:                  "CVM20201218P00003000",
			Key:                 "CVM",
			Subkey:              "CVM20201218P00003000",
			Date:                time.Date(2020, 11, 25, 0, 0, 0, 0, time.UTC),
			Updated:             time.Date(2020, 11, 30, 12, 25, 17, 0, time.UTC),
		}}
	})

	It("should parse option prices correctly", func() {
		var res []Option
		helper.TestdataFromJSON("core/option/options.json", &res)
		Expect(cmp.Equal(res, expected)).To(BeTrue(), cmp.Diff(res, expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the Option is valid", func() {
			Expect(expected[0].Validate()).To(Succeed())
		})
		It("should return an error if the Symbol is empty", func() {
			expected[0].Symbol = ""
			Expect(expected[0].Validate()).To(MatchError("missing symbol"))
		})
		It("should return an error if the ID is empty", func() {
			expected[0].ID = ""
			Expect(expected[0].Validate()).To(MatchError("missing id"))
		})
		It("should return an error if the ExpirationDate is zero valued", func() {
			expected[0].ExpirationDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("missing expiration date"))
		})
		It("should return an error if the StrikePrice is zero", func() {
			expected[0].StrikePrice = 0
			Expect(expected[0].Validate()).To(MatchError("strike price is zero"))
		})
	})

	It("should load the golden file", func() {
		// Pointers to zero values don't encode/decode properly
		expected[0].MarginPrice = nil
		expected[0].SettlementPrice = nil
		var golden []Option
		helper.FromGolden("option", &golden)
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("option", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})
})
