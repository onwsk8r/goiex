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

// +build !integration

package stock_test

import (
	"context"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
	"github.com/onwsk8r/goiex/pkg/core/stock/price"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Price", func() {
	var p *Price

	BeforeEach(func() {
		p = NewPrice(client)
	})

	Describe("HistoricalDaily", func() {
		var params *HistoricalPriceParams

		BeforeEach(func() {
			params = new(HistoricalPriceParams)
		})

		It("should get all the prices", func() {
			helper.TestdataResponder("/stable/stock/twtr/chart", "core/stock/price/historical.json")

			res, err := p.HistoricalDaily(context.Background(), "twtr", params)
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(ConsistOf(price.GoldenHistorical()))
		})

		It("should work with a dynamic range", func() {
			helper.TestdataResponder("/stable/stock/twtr/chart", "core/stock/price/historical_dynamic.json")

			params.Range = "dynamic"
			res, err := p.HistoricalDaily(context.Background(), "twtr", params)
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(ConsistOf(price.GoldenHistorical()))
		})
	})

	Describe("Intraday", func() {
		var params *HistoricalPriceParams

		BeforeEach(func() {
			params = new(HistoricalPriceParams)
			helper.TestdataResponder("/stable/stock/twtr/intraday-prices", "core/stock/price/intraday.json")
		})

		It("should get all the prices", func() {
			res, err := p.Intraday(context.Background(), "twtr", params)
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())

			expected := price.GoldenIntraday()
			Expect(res[0].Date.Equal(expected.Date)).To(BeTrue(), "dates are inequal")
			res[0].Date = expected.Date
			Expect(res[0]).To(Equal(expected))
		})
	})

	Describe("PreviousDay", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/twtr/previous", "core/stock/price/previous_day.json")
		})

		It("should get all the prices", func() {
			res, err := p.PreviousDay(context.Background(), "twtr")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(price.GoldenPreviousDay()))
		})
	})

	Describe("PreviousDayMarket", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/market/previous", "core/stock/price/previous_day_market.json")
		})

		It("should get all the prices", func() {
			res, err := p.PreviousDayMarket(context.Background())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(ConsistOf(price.GoldenPreviousDay()))
		})
	})
})
