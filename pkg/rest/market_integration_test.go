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

var _ = Describe("Market", func() {
	var m *Market

	BeforeEach(func() {
		m = NewMarket(client)
		Expect(m).ToNot(BeNil())
	})

	Describe("UpcomingDividends", func() {
		It("should successfully get and parse upcoming dividends", func() {
			res, err := m.UpcomingDividends(ctx, "market")
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			// Historically there have always been >500 dividends, but a build
			// failed on account of only finding 339. Regardless, expect a bunch.
			Expect(len(res)).To(BeNumerically(">", 100))
			Expect(len(res)).To(BeNumerically("<", 2500))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})

	Describe("UpcomingSplits", func() {
		It("should successfully get and parse upcoming splits", func() {
			res, err := m.UpcomingSplits(ctx, "market")
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically(">", 1))
			Expect(len(res)).To(BeNumerically("<", 1000))

			for idx := range res {
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})
})
