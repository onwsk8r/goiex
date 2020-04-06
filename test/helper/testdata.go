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

package helper

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/gomega" // nolint: stylecheck
)

// Testdata fetches a testdata file.
// The path should be relative to the /test/testdata directory.
func Testdata(path string) (io.ReadCloser, error) {
	path = filepath.Join(RepoBaseDir, "test", "testdata", filepath.FromSlash(path))
	return os.Open(path)
}

// TestdataFromJSON fetches a testdata file and decodes it into a type.
// The path should be compatible with Testdata(), and the target should
// be compatible with Golang's json.Unmarshal().
func TestdataFromJSON(path string, target interface{}) {
	rc, err := Testdata(path)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	defer rc.Close()

	dec := json.NewDecoder(rc)
	dec.DisallowUnknownFields()
	ExpectWithOffset(1, dec.Decode(target)).To(Succeed())
}

// TestdataReponder registers an httpmock responder that responds with the given testdata.
func TestdataResponder(url, testdata string) {
	rc, err := Testdata(testdata)
	ExpectWithOffset(1, err).ToNot(HaveOccurred(), "error loading testdata file")
	data, err := ioutil.ReadAll(rc)
	ExpectWithOffset(1, err).ToNot(HaveOccurred(), "error reading testdata file")
	httpmock.RegisterResponder("GET", url,
		httpmock.NewBytesResponder(http.StatusOK, data))
}
