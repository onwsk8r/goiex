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

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/stock"
	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Fundamentals", func() {
	var f *Fundamentals

	BeforeEach(func() {
		f = NewFundamentals(client)
	})

	Describe("Dividends", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/twtr/dividends/1y", "core/stock/fundamental/dividends.json")
		})

		It("should get all the dividends", func() {
			res, err := f.Dividends(context.Background(), "twtr", "1y")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			expected := fundamental.GoldenDividends()
			Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
		})
	})

	Describe("Earnings", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/twtr/earnings/2", "core/stock/fundamental/earnings.json")
		})

		It("should get all the earnings", func() {
			res, err := f.Earnings(context.Background(), "twtr", 2, false)
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			expected := fundamental.GoldenEarnings()
			Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
		})
	})

	Describe("Splits", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/twtr/splits/5y", "core/stock/fundamental/splits.json")
		})

		It("should get all the splits", func() {
			res, err := f.Splits(context.Background(), "twtr", "5y")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			expected := fundamental.GoldenSplit()
			Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
		})
	})
})
