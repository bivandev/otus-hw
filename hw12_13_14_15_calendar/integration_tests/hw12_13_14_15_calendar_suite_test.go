package integrationtests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHw12131415Calendar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hw12131415Calendar Suite")
}
