package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/bunniesandbeatings/commandgo"
	. "github.com/bunniesandbeatings/commandgo/ginkgocumber"
	. "github.com/MakeNowJust/heredoc/dot"
	"regexp"
)

// NOTE: There are no flags tests Here . Trusting that ExtractOptions is sufficient to specify behavior.

var _ = Describe("extract sub-command", func() {

	var (
		runner *Runner
	)

	Given("Statefile as stdin", func() {
		BeforeEach(func() {
			runner = NewRunner(offspringCLI, "extract")
		})

		stateString := D(`
			---
			a-password: 1234ILikethings
			fiddly-password: i-am-sensitive
			url: https://username:i-am-sensitive@basic-auth.com
		`)

		When("I dont match an expression", func() {
			BeforeEach(func() {
				runner.AddArguments("-k", "password")
				runner.AddArguments("-p", "i-dont-match")
			})

			Then("statefile is at stdout", func() {
				session := runner.ExecuteWithInput(stateString)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(stateString)))
				Eventually(session).Should(Exit(0))
			})
		})

		When("I search for my pattern expresssion", func() {
			BeforeEach(func() {
				runner.AddArguments("-k", "password")
				runner.AddArguments("-p", "1234ILikethings")
			})

			Then("statefile is at stdout with my credential templatized", func() {
				session := runner.ExecuteWithInput(stateString)

				templatized := D(`
					---
					a-password: {{password}}
					fiddly-password: i-am-sensitive
					url: https://username:i-am-sensitive@basic-auth.com
				`)

				Eventually(session.Out).Should(Say(regexp.QuoteMeta(templatized)))
				Eventually(session).Should(Exit(0))
			})

			Then("my credential is written to a file", func() {

			})
		})
	})

	//
	//var runner *Runner
	//var credentialFileName string
	//
	//BeforeEach(func() {
	//	credentialFile, _ := ioutil.TempFile("", "credential-file") // TODO: candidate for commandgo/outputfile
	//	credentialFile.Close()
	//
	//	credentialFileName = credentialFile.Name()
	//
	//	runner = NewRunner(offspringCLI, "extract")
	//	runner.AddArguments("-f", credentialFile.Name())
	//})
	//
	//Describe("Statefile as an argument", func() {
	//	BeforeEach(func() {
	//		statefileContents := D(`
	//			---
	//			too-easy: found-me
	//			also-too-easy: found-me
	//			name: sensitive-name
	//			usage-of-name: sensitive-name.sensitive-site.com
	//		`)
	//
	//		statefileFixture := NewFixture("state-file").
	//			Write([]byte(statefileContents)).
	//			Close()
	//
	//		runner.AddArguments("-s", statefileFixture.Name())
	//	})
	//
	//	Context("With a simple pattern that matches once", func() {
	//		BeforeEach(func() {
	//			runner.AddArguments("-k", "name-line")
	//			runner.AddArguments("-p", `(?m)^name: .*`)
	//		})
	//
	//		It("extracts only the sensitive data", func() {
	//			command := runner.Command()
	//
	//			session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//			Expect(err).ToNot(HaveOccurred())
	//
	//			Eventually(session).Should(Exit(0))
	//			Expect(session).To(Say(D(`
	//			---
	//			too-easy: found-me
	//			also-too-easy: found-me
	//			{{name-line}}
	//			usage-of-name: sensitive-name.sensitive-site.com
	//		`)))
	//		})
	//	})
	//
	//	Context("With a capture pattern that matches once", func() {
	//		BeforeEach(func() {
	//			runner.AddArguments("-k", "tld")
	//			runner.AddArguments("-p", `(?m)usage-of-name: [^\.]*\.(.*)$`)
	//		})
	//
	//		It("extracts only the sensitive data", func() {
	//			command := runner.Command()
	//
	//			session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//			Expect(err).ToNot(HaveOccurred())
	//
	//			Eventually(session).Should(Exit(0))
	//			Expect(session).To(Say(D(`
	//				---
	//				too-easy: found-me
	//				also-too-easy: found-me
	//				name: sensitive-name
	//				usage-of-name: sensitive-name.{{tld}}
	//			`)))
	//		})
	//	})
	//
	//	Context("With a simple pattern that matches more than once", func() {
	//		BeforeEach(func() {
	//			runner.AddArguments("-k", "name-line")
	//			runner.AddArguments("-p", `name: .*`)
	//		})
	//
	//		It("tells you it overmatched", func() {
	//			command := runner.Command()
	//
	//			session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//			Expect(err).ToNot(HaveOccurred())
	//
	//			Eventually(session).Should(Exit(1))
	//
	//			Expect(session.Err).To(Say(regexp.QuoteMeta(`Matched pattern 'name: .*' 2 time(s) in the state-file, you need to improve the match count`)))
	//		})
	//
	//		//Context("with the multiple flag", func() {
	//		//	BeforeEach(func() {
	//		//		runner.AddArguments("-m")
	//		//	})
	//		//
	//		//	It("extracts both occurances of the match", func() {
	//		//		command := runner.Command()
	//		//
	//		//		session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//		//		Expect(err).ToNot(HaveOccurred())
	//		//
	//		//		Eventually(session).Should(Exit(0))
	//		//		Expect(session).To(Say(D(`
	//		//			---
	//		//			too-easy: found-me
	//		//			also-too-easy: found-me
	//		//			{{name-line}}
	//		//			usage-of-{{name-line}}
	//		//		`)))
	//		//	})
	//		//})
	//	})
	//
	//	XContext("With a capture pattern that matches more than once", func() {
	//		BeforeEach(func() {
	//			runner.AddArguments("-k", "found")
	//			runner.AddArguments("-p", `(?ms)easy: (.*?)$`)
	//		})
	//
	//		It("tells you it overmatched the catpture group", func() {
	//			command := runner.Command()
	//
	//			session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//			Expect(err).ToNot(HaveOccurred())
	//
	//			Eventually(session).Should(Exit(1))
	//			Expect(session.Err).To(Say(regexp.QuoteMeta(`Matched group '(?s)easy: (.*?)' 2 time(s) in the state-file, you need to improve the match count or enable multiple matches with '-m'`)))
	//		})
	//
	//		Context("with the multiple flag", func() {
	//			BeforeEach(func() {
	//				runner.AddArguments("-m")
	//			})
	//
	//			It("extracts only the sensitive data", func() {
	//				command := runner.Command()
	//
	//				session, err := Start(command, GinkgoWriter, GinkgoWriter)
	//				Expect(err).ToNot(HaveOccurred())
	//
	//				Eventually(session).Should(Exit(0))
	//				Expect(session).To(Say(D(`
	//					---
	//					too-easy: {{found}}
	//					also-too-easy: {{found}}
	//					name: sensitive-name
	//					usage-of-name: sensitive-name.sensitive-site.com
	//				`)))
	//			})
	//		})
	//	})
	//
	//})
	//
	//Context("With a pattern that does not match", func() {
	//
	//})
	//
	//Context("With an invalid regex pattern", func() {
	//
	//})
})
