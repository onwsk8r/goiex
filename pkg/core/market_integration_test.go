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

package core_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core"
)

var _ = Describe("Market Integration", func() {
	var mkt *Market
	BeforeEach(func() {
		mkt = NewMarket(client)
	})

	Describe("UpcomingDividends", func() {
		It("should get dividends for the market", func() {
			res, err := mkt.UpcomingDividends(context.Background(), "market")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically(">", 1000))

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})

	Describe("UpcomingEarnings", func() {
		It("should get earnings for the market", func() {
			res, err := mkt.UpcomingEarnings(context.Background(), "market")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically(">", 1000))

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})

	Describe("UpcomingSplits", func() {
		It("should get splits for the market", func() {
			res, err := mkt.UpcomingEarnings(context.Background(), "market")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically(">", 500))

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})
})
