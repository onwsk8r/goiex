package stock_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPrice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stock Suite")
}
