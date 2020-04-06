package helper

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/gomega" // nolint: stylecheck
)

// ToGolden gob encodes val to a file called <name>.golden in the caller's directory
func ToGolden(name string, val interface{}) {
	_, file, _, ok := runtime.Caller(1)
	ExpectWithOffset(1, ok).To(BeTrue(), "error fetching runtime information")

	path := filepath.Join(filepath.Dir(file), fmt.Sprintf("%s.golden", name))
	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	ExpectWithOffset(1, err).ToNot(HaveOccurred(), fmt.Sprintf("error opening file %s", path))
	ExpectWithOffset(1, gob.NewEncoder(fh).Encode(val)).To(Succeed())
}

// FromGolden gob decodes <name>.golden in the caller's directory into val
func FromGolden(name string, val interface{}) {
	_, file, _, ok := runtime.Caller(1)
	ExpectWithOffset(1, ok).To(BeTrue(), "error fetching runtime information")

	path := filepath.Join(filepath.Dir(file), fmt.Sprintf("%s.golden", name))
	fh, err := os.Open(path)
	ExpectWithOffset(1, err).ToNot(HaveOccurred(), fmt.Sprintf("error opening file %s", path))
	ExpectWithOffset(1, gob.NewDecoder(fh).Decode(val)).To(Succeed())
}
