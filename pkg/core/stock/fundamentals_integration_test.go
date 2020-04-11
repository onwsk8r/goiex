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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
)

var _ = Describe("Fundamentals Integration", func() {
	var f *Fundamentals
	BeforeEach(func() {
		f = NewFundamentals(client)
	})

	Describe("Dividends", func() {
		It("should get the dividends", func() {
			res, err := f.Dividends(context.Background(), "ngl", "1y")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically("~", 4, 1))

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})

	Describe("Earnings", func() {
		It("should get the earnings", func() {
			res, err := f.Earnings(context.Background(), "aapl", 2, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(HaveLen(2))

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})

	Describe("Splits", func() {
		It("should get the splits", func() {
			res, err := f.Splits(context.Background(), "ko", "5y")
			Expect(err).ToNot(HaveOccurred())

			for _, val := range res {
				Expect(val.Validate()).To(Succeed())
			}
		})
	})
})
