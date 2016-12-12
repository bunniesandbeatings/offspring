package processor

import (
	"regexp"
	"fmt"
	"bytes"
	"errors"
)

type MultipleMatchError struct {
	Count int
}

func (error MultipleMatchError) Error() string {
	return fmt.Sprintf("Matched %s times without multiple mode", error.Count)
}

type Extraction struct {
	Name    string
	Pattern *regexp.Regexp // perhaps an re
	Multi   bool
}

type Extractor interface {
	Execute(source []byte) ([]byte, string, error)
}

func NewConfiguration(name string, pattern string, allowMultiple bool) (*Extraction, error) {
		patternExpression, regexpCompileError := regexp.Compile(pattern)
		if regexpCompileError != nil {
			return nil, regexpCompileError
		}

	return &Extraction{
		Name:    name,
		Pattern: patternExpression,
		Multi: allowMultiple,
	},
	nil
}

func (extraction *Extraction) Execute(source []byte) ([]byte, string, error) {

	found := extraction.Pattern.FindAllSubmatch(source, -1)

	//log.Printf("^^%s^^", found)

	if len(found) == 0 {
		return source, "", errors.New("Could not match expression")
	} else if (len(found) > 1) && !extraction.Multi {
		return source, "", MultipleMatchError{len(found)}
	}

	match := found[0]

	templateString := fmt.Sprintf("{{%s}}", extraction.Name)
	templateBytes := []byte(templateString)

	var password, out []byte

	if len(match) == 1 {
		password = match[0]
		out = extraction.Pattern.ReplaceAll(source, templateBytes)
	} else {
		password = match[1]
		out = bytes.Replace(source, password, templateBytes, 1)
	}

	return out, string(password), nil
}
