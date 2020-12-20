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

package core_test

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core"
	"github.com/onwsk8r/goiex/pkg/core/option"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Option", func() {
	var opt *Option

	BeforeEach(func() {
		opt = NewOption(client)
	})

	Describe("Options", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/stock/AAPL/options/201904", "core/option/options.json")
		})

		It("should get all the options", func() {
			expected := option.GoldenOption()
			res, err := opt.Options(context.Background(), "AAPL", "201904")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			// See the option entity class for why this is necessary
			res[0].SettlementPrice = nil
			res[0].MarginPrice = nil
			Expect(cmp.Equal(expected, res)).To(BeTrue(), cmp.Diff(expected, res))
		})
	})
})
