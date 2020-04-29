package option_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Option Suite")
}
