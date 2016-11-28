package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
)

type InjectOptions struct {
	Credential string `short:"c" long:"credential" description:"the credential to interpolate."`
	CredentialFile string `short:"f" long:"credential-file" description:"the credential to interpolate in a file"`

	StateFile string `short:"s" long:"state-file" description:"The state file. Defaults to STDIN"`

	Key string `short:"k" long:"key" description:"the moustache name for the credential to replace" required:"true"`
}

var injectOptions InjectOptions

func (injectOptions *InjectOptions) OpenStateFile() ([]byte, error) {
	if(injectOptions.StateFile == "") {
		output, stdoutError := ioutil.ReadAll(os.Stdin)
		if(stdoutError != nil) {
			return nil, errors.New("Could not read from stdin")
		}

		return output, nil

	}

	output, readFileError := ioutil.ReadFile(injectOptions.StateFile)
	if (readFileError != nil) {
		return nil, errors.New(fmt.Sprintf("Could not read input file '%s'", injectOptions.StateFile))
	}

	return output, nil
}

func (injectOptions *InjectOptions) Execute(args []string) error {
	state, stateFileError := injectOptions.OpenStateFile()

	if (stateFileError != nil) {
		return stateFileError
	}

	var credentials []byte

	if (injectOptions.CredentialFile != "") {
		var credentialsFileError error
		credentials, credentialsFileError = ioutil.ReadFile(injectOptions.CredentialFile)

		if (credentialsFileError != nil) {
			return errors.New(fmt.Sprintf("Could not open credentials file '%s'", injectOptions.CredentialFile))
		}

	} else if (injectOptions.Credential != "") {
		credentials = []byte(injectOptions.Credential)
	} else {
		return errors.New("You must provide either a credential string or a credential file")
	}


	// TODO: error out if template is not present

	handlebar, _ := regexp.Compile(regexp.QuoteMeta("{{" + injectOptions.Key + "}}"))

	newState := handlebar.ReplaceAll(state, credentials)

	fmt.Print(newState)

	return nil
}

func init() {
	parser.AddCommand(
		"inject",
		"Inject Credentials",
		"The inject command injects a credentials file or string into the state file",
		&injectOptions)
}
