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

package api_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/internal/api"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
)

func createClientAndEmptyRequest(token, version string) {
	client, err := NewClient(token, version)
	Expect(err).ToNot(HaveOccurred())

	err = client.Get(context.Background(), nil, url.Values{}, func(io.ReadCloser) error { return nil })
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("NewClient", func() {
	It("should use the stable version of the API by default", func() {
		token := "pk_sometoken"
		httpmock.RegisterResponder("GET", fmt.Sprintf("/%s?token=%s", iexcloud.APIVersionStable, token),
			httpmock.NewStringResponder(200, "OK"))
		createClientAndEmptyRequest(token, "")
	})
	It("should use the correct domain and version", func() {
		token := "pk_sometoken"
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s?token=%s",
			iexcloud.APIDomainBase, iexcloud.APIVersionBeta, token),
			httpmock.NewStringResponder(200, "OK"))
		createClientAndEmptyRequest(token, iexcloud.APIVersionBeta)
	})
	It("should recognize a sandbox token and use the sandbox domain", func() {
		By("recognizing test private keys")
		token := "Tpk_sometoken"
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s?token=%s",
			iexcloud.APIDomainSandbox, iexcloud.APIVersionStable, token),
			httpmock.NewStringResponder(200, "OK"))
		createClientAndEmptyRequest(token, "")

		By("recognizing test shareable keys")
		token = "Tsk_sometoken"
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s?token=%s",
			iexcloud.APIDomainSandbox, iexcloud.APIVersionStable, token),
			httpmock.NewStringResponder(200, "OK"))
		createClientAndEmptyRequest(token, "")
	})
	It("should panic if the length of the token is < 4", func() {
		dangerous := func() { NewClient("123", "stable") } // nolint: go-lint
		Expect(dangerous).To(Panic())
	})
})

var _ = Describe("Client", func() {
	var client iexcloud.Client

	BeforeEach(func() {
		var err error
		client, err = NewClient("pk_sometoken", "stable")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should not make a request if the context has been canceled", func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc()
		err := client.Get(ctx, nil, url.Values{}, func(io.ReadCloser) error { return nil })
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ctx.Err()))
	})

	Context("with a 200 return status", func() {
		var body string

		BeforeEach(func() {
			body = "beep boop"
			httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, body))
		})

		It("should execute the passed function", func() {
			err := client.Get(context.Background(), nil, url.Values{}, func(io.ReadCloser) error {
				Succeed()
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass the response body to the passed function", func() {
			err := client.Get(context.Background(), nil, url.Values{}, func(r io.ReadCloser) error {
				res, err := ioutil.ReadAll(r)
				Expect(err).ToNot(HaveOccurred())
				Expect(res).To(BeEquivalentTo(body))
				return err
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return errors from the passed function", func() {
			err := client.Get(context.Background(), nil, url.Values{}, func(io.ReadCloser) error {
				return fmt.Errorf("The %s was %s", "data", "frobulated")
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("The data was frobulated"))
		})
	})

	Context("with a non-200 return status", func() {
		var body string

		BeforeEach(func() {
			body = "Not Found"
			httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, body))
		})

		It("should not invoke the passed function", func() {
			err := client.Get(context.Background(), nil, url.Values{}, func(io.ReadCloser) error {
				Fail("passed function invoked")
				return nil
			})
			Expect(err).To(HaveOccurred())
		})

		It("should return an error with the HTTP status code and message", func() {
			err := client.Get(context.Background(), nil, url.Values{}, func(io.ReadCloser) error { return nil })
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("received HTTP status code 404 (404)"))
		})
	})
})
