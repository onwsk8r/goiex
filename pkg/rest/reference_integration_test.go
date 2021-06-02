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

var _ = Describe("Reference", func() {
	var r *Reference

	BeforeEach(func() {
		r = NewReference(client)
		Expect(r).ToNot(BeNil())
	})

	Describe("Symbols", func() {
		It("should successfully get and parse symbols", func() {
			res, err := r.Symbols(ctx)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically(">", 2500))
			Expect(len(res)).To(BeNumerically("<", 15000))

			for idx := range res {
				err = res[idx].Validate()
				Expect(err).To(SatisfyAny(BeNil(), MatchError("missing CIK")))
			}
		})
	})

	Describe("OptionsSymbols", func() {
		It("should successfully get and parse option symbols", func() {
			res, err := r.OptionsSymbols(ctx)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically(">", 2500))
			Expect(len(res)).To(BeNumerically("<", 10000))

			for idx := range res {
				Expect(len(idx)).To(SatisfyAll(BeNumerically(">=", 1), BeNumerically("<", 16)))
				Expect(len(res[idx])).To(SatisfyAll(BeNumerically(">=", 1), BeNumerically("<", 60)), idx)
			}
		})
	})
})
