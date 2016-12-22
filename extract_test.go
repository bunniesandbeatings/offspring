package main_test

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/bunniesandbeatings/commandgo"
	. "github.com/bunniesandbeatings/commandgo/ginkgocumber"
	. "github.com/MakeNowJust/heredoc/dot"
	"regexp"
	"io/ioutil"
)

// NOTE: There are no flags tests here.
// Trusting that ExtractOptions is sufficient to specify behavior.

var _ = Describe("extract sub-command", func() {

	var (
		executable         *ExecutableContext
		credentialFilePath string
		session            *Session
	)

	Given("Statefile as stdin", func() {
		BeforeEach(func() {
			disk := NewDisk()
			credentialFilePath = disk.CreateTempFilePath("credential-file-")

			executable = NewExecutableContext(offspringCLI, "extract")
			executable.AddArguments("-f", credentialFilePath)
		})

		stateString := D(`
			---
			a-password: 1234ILikethings
			fiddly-password: i-am-sensitive
			url: https://username:i-am-sensitive@basic-auth.com
		`)

		When("I dont match an expression", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "i-dont-match")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("statefile is at stdout", func() {
				Eventually(session.Out).Should(Say(regexp.QuoteMeta(stateString)))
				Eventually(session).Should(Exit(0))
			})
		})

		When("I search for my pattern expresssion", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "1234ILikethings")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("statefile is at stdout with my credential templatized", func() {
				templatized := D(`
					---
					a-password: {{password}}
					fiddly-password: i-am-sensitive
					url: https://username:i-am-sensitive@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))

				credentialFileContents, err := ioutil.ReadFile(credentialFilePath)

				Expect(err).To(BeNil())
				Expect(string(credentialFileContents)).To(Equal("1234ILikethings"))
			})
		})

		When("I explicitly want curly mustaches", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "1234ILikethings")
				executable.AddArguments("--mustache", "curly")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("statefile is at stdout with my credential templatized", func() {
				templatized := D(`
					---
					a-password: {{password}}
					fiddly-password: i-am-sensitive
					url: https://username:i-am-sensitive@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))
			})
		})

		When("I explicitly want smooth mustaches", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "1234ILikethings")
				executable.AddArguments("--mustache", "smooth")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("statefile is at stdout with my credential templatized all smooth like", func() {
				templatized := D(`
					---
					a-password: ((password))
					fiddly-password: i-am-sensitive
					url: https://username:i-am-sensitive@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))
			})
		})



		When("I search for a group capturing expresssion", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "a-password: (.*)")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("statefile is at stdout with my credential templatized", func() {
				templatized := D(`
					---
					a-password: {{password}}
					fiddly-password: i-am-sensitive
					url: https://username:i-am-sensitive@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))

				credentialFileContents, err := ioutil.ReadFile(credentialFilePath)

				Expect(err).To(BeNil())
				Expect(string(credentialFileContents)).To(Equal("1234ILikethings"))
			})
		})

		When("I unintentionally search for an expression that matches twice", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "i-am-sensitive")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("the user is informed that multi must be enables", func() {
				Eventually(session.Err).Should(Say("Pattern matches 2 times. Do you want the -m flag?"))
				Eventually(session).Should(Exit(1))

				credentialFileContents, err := ioutil.ReadFile(credentialFilePath)
				Expect(err).To(BeNil())
				Expect(string(credentialFileContents)).To(Equal(""))
			})
		})

		When("I intentionally search for an expression that matches twice", func() {
			BeforeEach(func() {
				executable.AddArguments("-k", "password")
				executable.AddArguments("-p", "i-am-sensitive")
				executable.AddArguments("-m")

				session = executable.ExecuteWithInput(stateString)
			})

			Then("the user is informed that multi must be enables", func() {
				templatized := D(`
					---
					a-password: 1234ILikethings
					fiddly-password: {{password}}
					url: https://username:{{password}}@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))

				credentialFileContents, err := ioutil.ReadFile(credentialFilePath)

				Expect(err).To(BeNil())
				Expect(string(credentialFileContents)).To(Equal("i-am-sensitive"))
			})
		})

		When("I unintentionally search for a group matcher that matches twice", func() {
		})

		When("I intentionally search for a group matcher that matches twice", func() {
		})

	})
})
