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
	"github.com/onwsk8r/goiex/pkg/core/market"
	"github.com/onwsk8r/goiex/pkg/core/stock"

	. "github.com/onwsk8r/goiex/pkg/rest"
)

var _ = Describe("Market", func() {
	var m *Market

	BeforeEach(func() {
		m = NewMarket(client)
		Expect(m).ToNot(BeNil())
	})

	Describe("UpcomingDividends", GetAndVerify("/v1/stock/NGL/upcoming-dividends", market.GoldenUpcomingDividends(),
		func() (interface{}, error) { return m.UpcomingDividends(ctx, "NGL") }))

	Describe("UpcomingSplits", GetAndVerify("/v1/stock/AAPL/upcoming-splits", stock.GoldenSplit(),
		func() (interface{}, error) { return m.UpcomingSplits(ctx, "AAPL") }))
})
