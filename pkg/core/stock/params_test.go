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

package stock_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
)

var _ = Describe("HistoricalPriceParams", func() {
	var p HistoricalPriceParams
	BeforeEach(func() { p = HistoricalPriceParams{} })

	It("should print an empty string with an empty struct", func() {
		Expect(p.Values().Encode()).To(Equal(""))
	})

	DescribeTable("Possible query string values",
		func(p HistoricalPriceParams, key, value string) {
			want := fmt.Sprintf("%s=%s", key, value)
			got := p.Values().Encode()
			Expect(got).To(Equal(want))
		},
		Entry("chartCloseOnly", HistoricalPriceParams{ChartCloseOnly: true}, "chartCloseOnly", "true"),
		Entry("chartByDay", HistoricalPriceParams{ChartByDay: true}, "chartByDay", "true"),
		Entry("chartSimplify", HistoricalPriceParams{ChartSimplify: true}, "chartSimplify", "true"),
		Entry("chartInterval", HistoricalPriceParams{ChartInterval: 5}, "chartInterval", "5"),
		Entry("changeFromClose", HistoricalPriceParams{ChangeFromClose: true}, "changeFromClose", "true"),
		Entry("chartLast", HistoricalPriceParams{ChartLast: 10}, "chartLast", "10"),
		Entry("range", HistoricalPriceParams{Range: "5y"}, "range", "5y"),
		Entry("exactDate", HistoricalPriceParams{ExactDate: time.Now()}, "exactDate", time.Now().Format("20060102")),
		Entry("sort", HistoricalPriceParams{Sort: "asc"}, "sort", "asc"),
		Entry("includeToday", HistoricalPriceParams{IncludeToday: true}, "includeToday", "true"),
	)
})
