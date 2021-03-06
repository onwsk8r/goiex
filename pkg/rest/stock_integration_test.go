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

package rest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/rest"
)

var _ = Describe("Stock", func() {
	var s *Stock

	BeforeEach(func() {
		s = NewStock(client)
		Expect(s).ToNot(BeNil())
	})

	Describe("Dividends", func() {
		It("should successfully get and parse dividends", func() {
			res, err := s.Dividends(ctx, "GE", DividendsPeriod1y)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("~", 4, 1))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})

	Describe("Earnings", func() {
		It("should successfully get and parse earnings", func() {
			res, err := s.Earnings(ctx, "GOOG", map[string]string{"last": "4"})
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("==", 4))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})

	Describe("Historical", func() {
		It("should successfully get and parse historical prices", func() {
			res, err := s.Historical(ctx, "GOOG", HistoricalPeriod1y, nil)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("~", 260, 10))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})

	Describe("HistoricalIntraday", func() {
		It("should successfully get and parse intraday prices", func() {
			res, err := s.HistoricalIntraday(ctx, "GOOG", HistoricalIntradayPeriod5dm, nil)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("~", 160, 50))

			for idx := range res {
				Expect(res[idx].Validate()).To(
					SatisfyAny(Succeed(), MatchError("market close is zero")))
			}
		})
	})

	Describe("PreviousDay", func() {
		It("should successfully get and parse previous day prices", func() {
			res, err := s.PreviousDay(ctx, "IBM")
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(res.Validate()).To(Succeed())
		})
	})

	Describe("PreviousDayMarket", func() {
		It("should successfully get and parse previous day prices", func() {
			res, err := s.PreviousDayMarket(ctx)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("~", 9000, 3000))

			closeIsZeroCount := 0
			for idx := range res {
				Expect(res[idx].Validate()).To(Or(Succeed(), MatchError("close is zero")))
				// This is kind of a hackish way to do this - quick and dirty - and a more robust
				// solution would involve figuring out _when_ the close is zero and why: is there
				// a reason for it? Is this a data quality error from IEX? Etc.
				if err := res[idx].Validate(); err != nil && err.Error() == "close is zero" {
					closeIsZeroCount++
				}
			}
			Expect(closeIsZeroCount).To(BeNumerically("<", 50)) // 50/6-12k isn't bad, riiiight?
		})
	})

	Describe("Splits", func() {
		It("should successfully get and parse upcoming splits", func() {
			res, err := s.Splits(ctx, "AAPL", SplitsPeriod5y)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically("==", 1))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})
})
