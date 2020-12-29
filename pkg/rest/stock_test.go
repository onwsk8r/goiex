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

package rest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onwsk8r/goiex/pkg/core/stock"

	. "github.com/onwsk8r/goiex/pkg/rest"
)

var _ = Describe("Stock", func() {
	var s *Stock

	BeforeEach(func() {
		s = NewStock(client)
		Expect(s).ToNot(BeNil())
	})

	Describe("Dividends (basic)", GetAndVerify("/v1/stock/NGL/dividends/ytd", stock.GoldenDividends(),
		func() (interface{}, error) { return s.Dividends(ctx, "NGL", DividendsPeriodYTD) }))

	Describe("Earnings", func() {
		var expected []stock.Earning
		BeforeEach(func() { expected = stock.GoldenEarnings() })

		Context("With default settings", GetAndVerify("/v1/stock/AAPL/earnings", expected,
			func() (interface{}, error) { return s.Earnings(ctx, "AAPL", nil) }))
		Context("When a period is passed", GetAndVerify("/v1/stock/MSFT/earnings?period=annual&token=sk_sometoken", expected,
			func() (interface{}, error) { return s.Earnings(ctx, "MSFT", map[string]string{"period": "annual"}) }))
		Context("When the 'last' parameter is specified",
			GetAndVerify("/v1/stock/TWTR/earnings?last=4&token=sk_sometoken", expected,
				func() (interface{}, error) { return s.Earnings(ctx, "TWTR", map[string]string{"last": "4"}) }))
	})

	Describe("Historical", func() {
		var expected []stock.Historical
		BeforeEach(func() { expected = stock.GoldenHistorical() })

		Context("With default settings", GetAndVerify("/v1/stock/TWTR/chart/ytd", expected,
			func() (interface{}, error) { return s.Historical(ctx, "TWTR", HistoricalPeriodYTD, nil) }))
		Context("With query string parameters",
			GetAndVerify("/v1/stock/MSFT/chart/ytd?chartLast=10&token=sk_sometoken", expected,
				func() (interface{}, error) {
					params := map[string]string{"chartLast": "10"}
					return s.Historical(ctx, "MSFT", HistoricalPeriodYTD, params)
				}))
	})

	Describe("HistoricalIntraday", func() {
		var expected []stock.Intraday
		BeforeEach(func() { expected = stock.GoldenIntraday() })

		Context("With only a range", GetAndVerify("/v1/stock/TWTR/chart/1mm", expected,
			func() (interface{}, error) {
				return s.HistoricalIntraday(ctx, "TWTR", HistoricalIntradayPeriod1mm, nil)
			}))
		Context("With a specific date",
			GetAndVerify("/v1/stock/MSFT/chart/date?exactDate=20190220&token=sk_sometoken", expected,
				func() (interface{}, error) {
					params := map[string]string{"exactDate": "20190220"}
					return s.HistoricalIntraday(ctx, "MSFT", HistoricalIntradayPeriodDate, params)
				}))
	})

	Describe("PreviousDay", GetAndVerify("/v1/stock/CSCO/previous", &stock.GoldenHistorical()[0],
		func() (interface{}, error) { return s.PreviousDay(ctx, "CSCO") }))

	Describe("PreviousDayMarket", GetAndVerify("/v1/stock/market/previous", stock.GoldenHistorical(),
		func() (interface{}, error) { return s.PreviousDayMarket(ctx) }))

	Describe("Splits (basic)", GetAndVerify("/v1/stock/AAPL/splits/5y", stock.GoldenSplit(),
		func() (interface{}, error) { return s.Splits(ctx, "AAPL", SplitsPeriod5y) }))
})
