package main_test

import (
	. "github.com/bunniesandbeatings/commandgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
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

		Context("Passing in credentials as a file", func() {
			BeforeEach(func() {
				credentialFixture := NewFixture("credential-file").
					Write([]byte(`VeryObscurePassword`)).
					Close()

				runner.AddArguments("-f", credentialFixture.Name())
			})

			Context("Passing in a matching credential key", func() {
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

			Context("Passing in a non-matching credential key", func() {
				// TODO: Perhaps a warn-if-not-matched method?
				BeforeEach(func() {
					runner.AddArguments("-k", "other-password")
				})

				It("leaves the output unchanged for pipeline based processing", func() {
					command, stdin := runner.PipeCommand()

					session, err := Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())

					stdin.Write([]byte("this is a {{sensitive-password}}"))
					stdin.Close()

					Eventually(session).Should(Say("this is a {{sensitive-password}}"))
					Eventually(session).Should(Exit(0))
				})
			})
		})
	})

	Describe("Statefile as a parameter", func() {
		BeforeEach(func() {
			statefileFixture := NewFixture("state-file").
				Write([]byte("this is a {{sensitive-password}}")).
				Close()

			runner.AddArguments("-s", statefileFixture.Name())
			runner.AddArguments("-c", "VeryObscurePassword")
			runner.AddArguments("-k", "sensitive-password")
		})

		It("adds the credential to the output", func() {
			command := runner.Command()

			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Expect(session).To(Say("this is a VeryObscurePassword"))
		})
	})
})
