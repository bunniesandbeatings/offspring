package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/bunniesandbeatings/commandgo/ginkgocumber"
	. "github.com/onsi/gomega"
	. "github.com/bunniesandbeatings/offspring/processor"
	. "github.com/MakeNowJust/heredoc/dot"

	"regexp/syntax"
)

var _ = Describe("NewConfiguration", func() {
	var (
		configuration *Extraction
		err           error
	)

	It("Is an Extractor", func() {
		var concrete interface{} = &Extraction{}
		_, ok := concrete.(Extractor)
		Expect(ok).To(BeTrue())
	})

	When("I construct an extraction configuration with a bad pattern", func() {
		BeforeEach(func() {
			configuration, err = NewConfiguration(
				"anything",
				"this(is(a(bad(regexp",
			)
		})

		Then("It returns a regexp compile error", func() {
			Expect(err).To(BeAssignableToTypeOf(&syntax.Error{}))
		})
	})

	When("I construct an extraction configuration with a good pattern", func() {
		BeforeEach(func() {
			configuration, err = NewConfiguration(
				"anything",
				"foo(.*)",
			)
			configuration.Multi = true
		})

		Then("It is built correctly", func() {
			Expect(err).To(BeNil())
			Expect(configuration.Multi).To(BeTrue())
			Expect(configuration.Name).To(Equal("anything"))
			Expect(configuration.Pattern.String()).To(Equal("foo(.*)"))
		})
	})
})

var _ = Describe("Execute", func() {
	var (
		newState      []byte
		credential    string
		err           error
		configuration *Extraction
	)

	Given("A single line with my state", func() {
		state := []byte("email: fred@example.com, password: furryknuckleduster")

		When("I extract the password explicitly", func() {
			BeforeEach(func() {
				configuration, _ = NewConfiguration("password", "furryknuckleduster")
			})

			Then("It extracts the password and tempaltizes it", func() {
				newState, credential, err = configuration.Execute(state)
				Expect(err).To(BeNil())
				Expect(string(newState)).To(Equal("email: fred@example.com, password: {{password}}"))
				Expect(credential).To(Equal("furryknuckleduster"))
			})

			When("I want the braces to be different", func() {
				BeforeEach(func() {
					configuration.LeftTash = "LLL"
				  configuration.RightTash = "RRR"
				})

				Then("It extracts the password and tempaltizes it", func() {
					newState, credential, err = configuration.Execute(state)
					Expect(err).To(BeNil())
					Expect(string(newState)).To(Equal("email: fred@example.com, password: LLLpasswordRRR"))
					Expect(credential).To(Equal("furryknuckleduster"))
				})
			})
		})

		When("I extract the password with a singular matching regex", func() {
			BeforeEach(func() {
				configuration, _ = NewConfiguration("password", "furry.*duster")
				newState, credential, err = configuration.Execute(state)
			})

			Then("It extracts the password and tempaltizes it", func() {
				Expect(err).To(BeNil())
				Expect(string(newState)).To(Equal("email: fred@example.com, password: {{password}}"))
				Expect(credential).To(Equal("furryknuckleduster"))
			})
		})

		When("I extract the password with a capture group regex", func() {
			BeforeEach(func() {
				configuration, _ = NewConfiguration("password", "password: (.*)")
				newState, credential, err = configuration.Execute(state)
			})

			Then("It extracts the password and tempaltizes it", func() {
				Expect(err).To(BeNil())
				Expect(string(newState)).To(Equal("email: fred@example.com, password: {{password}}"))
				Expect(credential).To(Equal("furryknuckleduster"))
			})
		})

		When("I use a regex that does not match", func() {
			BeforeEach(func() {
				configuration, _ = NewConfiguration("password", "IDontMatch")
				newState, credential, err = configuration.Execute(state)
			})

			Then("It returns a no-match error", func() {
				Expect(err).To(MatchError("Could not match expression"))
				Expect(string(newState)).To(Equal("email: fred@example.com, password: furryknuckleduster"))
				Expect(credential).To(Equal(""))
			})
		})
	})

	Given("A document with multiple matching opportunities", func() {
		state := []byte(D(`
			---
			a-password: 1234ILikethings
			fiddly-password: i-am-sensitive
			url: https://username:i-am-sensitive@basic-auth.com
		`))

		When("I match more than once", func() {
			BeforeEach(func() {
				configuration, _ = NewConfiguration("multi-pass", "i-am-sensitive")
			})

			Then("It returns a too many matches error", func() {
				newState, credential, err = configuration.Execute(state)
				Expect(err).To(MatchError(MultipleMatchError{Count: 2}))
				Expect(newState).To(Equal(state))
				Expect(credential).To(Equal(""))
			})

			When("I set multi to true", func() {
				BeforeEach(func() {
					configuration.Multi = true
				})

				Then("It templatizes all matches", func() {
					newState, credential, err = configuration.Execute(state)

					Expect(err).To(BeNil())
					Expect(string(newState)).To(Equal(D(`
						---
						a-password: 1234ILikethings
						fiddly-password: {{multi-pass}}
						url: https://username:{{multi-pass}}@basic-auth.com
					`)))
					Expect(credential).To(Equal("i-am-sensitive"))
				})

			})

		})

	})
})
