package speedmap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpeedmap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Speedmap Suite")
}
