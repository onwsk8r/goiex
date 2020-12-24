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
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	. "github.com/onwsk8r/goiex/pkg/rest"
)

var _ = Describe("REST Client", func() {
	Describe("The request URL", func() {
		It("should use a production domain", func() {
			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/foo", APIDomainBase),
				httpmock.NewStringResponder(http.StatusOK, "hello"))
			_, err := client.R().Get("/foo")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
		})
		It("should append the token as a QSP", func() {
			httpmock.RegisterResponder("GET", "/foo?token=sk_sometoken",
				httpmock.NewStringResponder(http.StatusOK, "hello"))
			_, err := client.R().Get("/foo")
			Expect(httpmock.GetTotalCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
		})
		Context("When using a sandbox token", func() {
			var sbClient *resty.Client
			JustBeforeEach(func() {
				sbClient = NewClient("Tsk_sometoken", log.Ctx(ctx))
				httpmock.ActivateNonDefault(sbClient.GetClient())
			})
			It("should use a sandbox domain", func() {
				httpmock.RegisterResponder("GET", fmt.Sprintf("%s/foo", APIDomainSandbox),
					httpmock.NewStringResponder(http.StatusOK, "hello"))
				_, err := sbClient.R().Get("/foo")
				Expect(httpmock.GetTotalCallCount()).To(Equal(1))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	It("should set a {version} path param", func() {
		httpmock.RegisterResponder("GET", fmt.Sprintf("/%s/foo", APIVersion1),
			httpmock.NewStringResponder(http.StatusOK, "hello"))
		_, err := client.R().Get("/{version}/foo")
		Expect(httpmock.GetTotalCallCount()).To(Equal(1))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with a canceled context", func() {
		var err error
		BeforeEach(func() {
			httpmock.RegisterNoResponder(httpmock.NewStringResponder(http.StatusOK, "hello"))
			newCtx, cancel := context.WithCancel(ctx)
			cancel()
			_, err = client.R().SetContext(newCtx).Get("/foo")
		})
		It("should short-circuit the request", func() { Expect(httpmock.GetTotalCallCount()).To(Equal(0)) })
		It("should return a 'context canceled' error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("context canceled"))
		})
	})

	Context("when an HTTP error occurs", func() {
		var err error
		BeforeEach(func() {
			httpmock.RegisterNoResponder(httpmock.NewErrorResponder(fmt.Errorf("oops")))
			_, err = client.R().Get("/foo")
		})
		It("should retry the request", func() { Expect(httpmock.GetTotalCallCount()).To(Equal(MaxRetries + 1)) })
		It("should return a *url.Error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err).To(BeAssignableToTypeOf(&url.Error{}))
		})
	})

	DescribeTable("Invalid response codes",
		func(code int, retry bool) {
			httpmock.RegisterNoResponder(httpmock.NewStringResponder(code, "hello"))
			_, err := client.R().Get("/foo")
			callCount := 1
			if retry {
				callCount = MaxRetries + 1
			}
			Expect(httpmock.GetTotalCallCount()).To(Equal(callCount))
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError(fmt.Sprintf("invalid response: %d hello", code)))
		},
		Entry("400 should not be retried", 400, false),
		Entry("401 should not be retried", 401, false),
		Entry("402 should not be retried", 402, false),
		Entry("403 should not be retried", 403, false),
		Entry("404 should not be retried", 404, false),
		Entry("413 should not be retried", 413, false),
		Entry("451 should not be retried", 451, false),
		Entry("429 should be retried", 429, true),
		Entry("500 should be retried", 500, true),
	)
})
