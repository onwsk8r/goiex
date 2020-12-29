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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/rest"
)

var _ = Describe("Options", func() {
	var o *Options

	BeforeEach(func() {
		o = NewOptions(client)
		Expect(o).ToNot(BeNil())
	})

	Describe("Expiration", func() {
		It("should successfully get and parse option expirations", func() {
			res, err := o.Expiration(ctx, "AAPL")
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("%+v", err))
			Expect(len(res)).To(BeNumerically(">", 10))

			now := time.Now().Add(-32 * 24 * time.Hour)
			later := now.Add(5 * 365 * 24 * time.Hour)
			var then time.Time
			for idx := range res {
				By(fmt.Sprintf("Having the correct value at %d (%s)", idx, res[idx]))
				Expect(res[idx]).To(SatisfyAny(HaveLen(6), HaveLen(8)),
					fmt.Sprintf("Invalid length (%d): '%s'", len(res[idx]), res[idx]))
				if len(res[idx]) == 6 {
					then, err = time.Parse("200601", res[idx])
				} else {
					then, err = time.Parse("20060102", res[idx])
				}
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Error parsing %s", res[idx]))
				Expect(then).To(SatisfyAll(BeTemporally(">", now), BeTemporally("<", later)))
			}
		})
	})

	Describe("EndOfDay", func() {
		It("should successfully get and parse options", func() {
			// Get next week's options
			date := time.Now().Add(7 * 24 * time.Hour)
			for date.Weekday() != time.Friday {
				date = date.Add(24 * time.Hour)
			}
			res, err := o.EndOfDay(ctx, "IBM", date.Format("20060102"))
			Expect(err).ToNot(HaveOccurred())
			if len(res) == 0 {
				// Options expire on Thursday if Friday is a holiday
				res, err = o.EndOfDay(ctx, "IBM", date.Add(-24*time.Hour).Format("20060102"))
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(len(res)).To(BeNumerically(">", 20))
			Expect(len(res)).To(BeNumerically("<", 60))
			for idx := range res {
				By(fmt.Sprintf("Having valid values at %d", idx))
				Expect(res[idx].Validate()).To(Succeed())
			}
		})
	})
})
