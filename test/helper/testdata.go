package helper

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
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
	Expect(err).ToNot(HaveOccurred())
	defer rc.Close()

	dec := json.NewDecoder(rc)
	dec.DisallowUnknownFields()
	Expect(dec.Decode(target)).To(Succeed())
}
