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

package rest_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onwsk8r/goiex/pkg/core/option"
	. "github.com/onwsk8r/goiex/pkg/rest"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Options", func() {
	var o *Options

	BeforeEach(func() {
		o = NewOptions(client)
		Expect(o).ToNot(BeNil())
	})

	GetAndVerify := func(url string, expected interface{}, f func() (interface{}, error)) func() {
		return func() {
			var got interface{}
			var err error
			BeforeEach(func() {
				httpmock.RegisterResponder("GET", url, httpmock.NewJsonResponderOrPanic(http.StatusOK, &expected))
				got, err = f()
			})
			It("should hit the expected URL", func() { Expect(httpmock.GetTotalCallCount()).To(Equal(1)) })
			It("should not encounter any errors", func() { Expect(err).ToNot(HaveOccurred()) })
			It("should return the expected data", func() {
				Expect(cmp.Equal(expected, got)).To(BeTrue(), cmp.Diff(expected, got))
			})
		}
	}

	Describe("Expiration", func() {
		var expected []string
		BeforeEach(func() { helper.TestdataFromJSON("core/option/expirations.json", &expected) })
		GetAndVerify("/v1/stock/AAPL/options", expected, func() (interface{}, error) {
			return o.Expiration(ctx, "AAPL")
		})()
	})

	Describe("EndOfDay", func() {
		Context("when getting a single side", GetAndVerify("/v1/stock/AAPL/options/20201231/call",
			option.GoldenOption(), func() (interface{}, error) {
				return o.EndOfDay(ctx, "AAPL", "20201231", "call")
			}))
		Context("when getting all options", GetAndVerify("/v1/stock/AAPL/options/20201231/",
			option.GoldenOption(), func() (interface{}, error) {
				return o.EndOfDay(ctx, "AAPL", "20201231")
			}))

	})
})
