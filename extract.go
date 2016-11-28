package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"bytes"
)

type ExtractOptions struct {
	StateFile string `short:"s" long:"state-file" description:"The state file. Defaults to STDIN"`

	Key string `short:"k" long:"key" description:"the moustache name for the credential placeholder" required:"true"`

	Pattern string `short:"p" long:"pattern" description:"The regular expression to find and extract the credentials (first capture group or entire expression)" required:"true"`
}

var extractOptions ExtractOptions

func (extractOptions *ExtractOptions) OpenStateFile() ([]byte, error) {
	if(extractOptions.StateFile == "") {
		output, stdoutError := ioutil.ReadAll(os.Stdin)
		if(stdoutError != nil) {
			return nil, errors.New("Could not read from stdin")
		}

		return output, nil
	}

	output, readFileError := ioutil.ReadFile(extractOptions.StateFile)
	if (readFileError != nil) {
		return nil, errors.New(fmt.Sprintf("Could not read input file '%s'", extractOptions.StateFile))
	}

	return output, nil
}

func (extractOptions *ExtractOptions) Execute(args []string) error {
	state, stateFileError := extractOptions.OpenStateFile()

	if (stateFileError != nil) {
		return stateFileError
	}

	pattern, regexpError := regexp.Compile(extractOptions.Pattern)

	if (regexpError != nil) {
		return errors.New(fmt.Sprintf("Could not compile pattern as regular expression '%s'", extractOptions.Pattern))
	}

	matches := pattern.FindAllSubmatch(state, -1)

	template := []byte("{{" + extractOptions.Key + "}}")

	var newState string

	if (len(matches) < 1) {
		return errors.New(fmt.Sprintf("Could not match pattern '%s' in the state-file", extractOptions.Pattern))
	} else if(len(matches) > 1) {
		return errors.New(fmt.Sprintf("Matched pattern '%s' %d times(s) in the state-file, you need to improve the match count", extractOptions.Pattern, len(matches)))
	} else {

		if (len(matches[0]) > 2) {
			return errors.New(fmt.Sprintf("Matched capture group %d times(s) in the state-file, cannot have more than one sub group match", len(matches[0])))
		} else if (len(matches[0]) > 1) {
			if (globalOptions.Debug) {
				fmt.Printf("Capture group match:\n\n%s",string(matches[0][1]))
				return nil
			}

			newState = string(pattern.ReplaceAllFunc(state, func(match []byte) []byte {
				return bytes.Replace(match, matches[0][1], template, 1)
			}))
		} else {
			if (globalOptions.Debug) {
				fmt.Printf("Whole expression match:\n\n%s",string(matches[0][0]))
				return nil
			}
			newState = string(pattern.ReplaceAll(state, template))
		}
	}

	fmt.Print(newState)
	return nil
}

func init() {
	parser.AddCommand(
		"extract",
		"Extract Credentials",
		"The extract command extracts credentials matching a regex and leaves a placeholder in the statefile",
		&extractOptions)
}
