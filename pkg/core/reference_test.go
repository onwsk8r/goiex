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

package core_test

import (
	"context"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("Reference", func() {
	var ref *Reference

	BeforeEach(func() {
		httpmock.Activate()
		ref = NewReference(client)
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	Describe("Symbols", func() {
		BeforeEach(func() {
			helper.TestdataResponder("/stable/ref-data/symbols", "core/reference/symbols.json")
		})

		It("should get all the symbols", func() {
			_, err := ref.Symbols(context.Background())
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
			// Expect(res).To(ConsistOf(expected))
		})
	})
})
