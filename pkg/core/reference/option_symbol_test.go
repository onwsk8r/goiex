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
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/reference"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("OptionSymbol", func() {
	var expected OptionSymbol

	BeforeEach(func() {
		expected = make(OptionSymbol)
		expected["DAL"] = []string{"201904", "201905", "201906", "201909", "201912", "202001", "202101"}
		expected["RUSS"] = []string{"201905", "201906", "201909", "201912"}
	})

	It("should parse symbols correctly", func() {
		var symbols OptionSymbol
		helper.TestdataFromJSON("core/reference/option_symbols.json", &symbols)
		Expect(symbols).To(BeEquivalentTo(expected))
	})

	It("should match the golden file", func() {
		golden := GoldenOptionSymbol()
		if !cmp.Equal(golden, expected) {
			helper.ToGolden("option_symbol", expected)
			Fail(cmp.Diff(golden, expected))
		}
	})
})
