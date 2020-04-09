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

// +build integration

package stock_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
)

var _ = Describe("Price Integration", func() {
	var p *Price
	BeforeEach(func() {
		p = NewPrice(client)
	})

	Describe("HistoricalDaily", func() {
		var params *HistoricalPriceParams

		BeforeEach(func() {
			params = new(HistoricalPriceParams)
		})

		It("should get one month of prices by default", func() {
			res, err := p.HistoricalDaily(context.Background(), "twtr", params)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically(">", 15))
			Expect(len(res)).To(BeNumerically("<", 50))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})

	Describe("Intraday", func() {
		var now time.Time
		var diff time.Duration
		BeforeEach(func() {
			// While there are 390 (or 210) minutes in a trading day, when called during
			// the trading day the result will have data for each minute from 9:30 - now.
			easternTime, err := time.LoadLocation("America/New_York")
			Expect(err).ToNot(HaveOccurred())
			now = time.Now().In(easternTime)
			open := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, easternTime)
			diff = now.Sub(open)
		})

		It("should get intraday prices", func() {
			res, err := p.Intraday(context.Background(), "twtr")
			Expect(err).ToNot(HaveOccurred())

			// There should be 390 or 210 results outside of normal trading hours
			// and just about "the number of minutes since 9:30" results during trading
			// hours, unless it's not a trading day (which we don't know)
			if diff < 0 || diff > time.Minute*390 {
				Expect(res).To(Or(HaveLen(390), HaveLen(210)))
			} else {
				// Give it a five-minute buffer, and account for non-trading days
				Expect(len(res)).To(Or(
					BeNumerically("~", diff.Minutes(), 5),
					BeNumerically("=", 390),
					BeNumerically("=", 210)))
			}

			for idx := range res {
				// Market close data is delayed fifteen minutes, and again with the 5-minute buffer.
				// The Close, which is the IEX close, may be zero if it did not trade on the IEX
				if now.Sub(res[idx].Date) > 20*time.Minute {
					Expect(res[idx].Validate()).To(Or(Succeed(), MatchError("close is zero")))
				} else {
					Expect(res[idx].Validate()).To(Or(Succeed(), MatchError("market close is zero")))
				}
			}
		})
	})

	Describe("PreviousDay", func() {
		It("should get the previous day's price for the given security", func() {
			res, err := p.PreviousDay(context.Background(), "twtr")
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Validate()).To(Succeed())
		})
	})

	Describe("PreviousDayMarket", func() {
		It("should get the previous day's prices for all securities", func() {
			res, err := p.PreviousDayMarket(context.Background())
			Expect(err).ToNot(HaveOccurred())

			// Some results will be empty and will not validate :(
			var valid int = 0
			for idx := range res {
				if err := res[idx].Validate(); err == nil {
					valid++
				}
			}
			// Hopefully no more than 50 are empty...
			Expect(valid).To(BeNumerically("~", len(res), 50))
		})
	})
})
