package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/bunniesandbeatings/commando/runner"
)


var _ = Describe("inject command", func() {

	var runner *Runner

	BeforeEach(func() {
		runner = NewRunner(offspringCLI,"inject")
	})

	Describe("Statefile as a pipe", func() {
		Describe("Passing in credentials as an argument", func() {
		  BeforeEach(func() {
				runner.AddArguments("-c", "VeryObscurePassword")
			})

			It("adds the credential to the output", func() {
				command, stdin := runner.PipeCommand("-k", "sensitive-password")

				session, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				stdin.Write([]byte("this is a {{sensitive-password}}"))
				stdin.Close()

				Eventually(session).Should(Say("this is a VeryObscurePassword"))
				Eventually(session).Should(Exit(0))
			})

		})

	})
})


