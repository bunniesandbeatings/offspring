package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/bunniesandbeatings/commandgo/runner"
)

var _ = Describe("inject command", func() {

	var runner *Runner

	BeforeEach(func() {
		runner = NewRunner(offspringCLI, "inject")
	})

	Describe("Statefile as a pipe", func() {
		Context("Passing in credentials as an argument", func() {
			BeforeEach(func() {
				runner.AddArguments("-c", "VeryObscurePassword")
			})

			Context("Passing in valid credential key", func() {
				BeforeEach(func() {
					runner.AddArguments("-k", "sensitive-password")
				})

				It("adds the credential to the output", func() {
					command, stdin := runner.PipeCommand()

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
})


