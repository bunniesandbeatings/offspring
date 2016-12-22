package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/bunniesandbeatings/offspring/processor"
	"errors"
)

type ExtractOptions struct {
	CredentialFile string `short:"f" long:"credential-file" description:"the credential file to create or update with the new creds" required:"true"`
	Key            string `short:"k" long:"key" description:"the template name for the credential placeholder" required:"true"`
	Pattern        string `short:"p" long:"pattern" description:"The regular expression to find and extract the credentials (first capture group or entire expression)" required:"true"`
	Multi bool `short:"m" long:"multi" description:"Replace multiple entries that match the capture group (will NOT work without capture group)"`
	//	StateFile string `short:"s" long:"state-file" description:"The state file. Defaults to STDIN"`

	// TODO DRY
	Mustache string `long:"mustache" description:"template format" choice:"curly" choice:"smooth" default:"curly"`

}

var extractOptions ExtractOptions

func (extractOptions *ExtractOptions) Execute(args []string) error {
	extractProcessor, _ := processor.NewConfiguration(extractOptions.Key, extractOptions.Pattern)
	extractProcessor.Multi = extractOptions.Multi
	credentialFile, _ := os.Create(extractOptions.CredentialFile)

	source, _ := ioutil.ReadAll(os.Stdin)

	if extractOptions.Mustache == "smooth" {
		extractProcessor.LeftTash = "(("
		extractProcessor.RightTash = "))"
	}

	output, credential, extractionError := extractProcessor.Execute(source)
	if matchErr, ok := extractionError.(processor.MultipleMatchError); ok {
		return errors.New(fmt.Sprintf("Pattern matches %d times. Do you want the -m flag?", matchErr.Count))
	}

	fmt.Fprint(credentialFile, credential)
	credentialFile.Close()

	fmt.Print(string(output))

	return nil
}

func init() {
	parser.AddCommand(
		"extract",
		"Extract Credentials",
		"The extract command extracts credentials matching a regex and leaves a placeholder in the statefile",
		&extractOptions)
}
