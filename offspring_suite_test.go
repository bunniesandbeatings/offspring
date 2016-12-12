package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/onsi/gomega/gexec"
)

func TestOffspring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offspring Suite")
}

var offspringCLI string

var _ = BeforeSuite(func() {
	var err error
	offspringCLI, err = gexec.Build("github.com/bunniesandbeatings/offspring")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})


