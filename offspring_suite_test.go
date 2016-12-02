package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "os/exec"

	"testing"
	"github.com/onsi/gomega/gexec"
	"log"
	"io"
)

func TestOffspring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offspring Suite")
}

var offspringCLI string

var _ = BeforeSuite(func() {
	var err error
	offspringCLI, err = gexec.Build("github.com/bunniesandbeatings/offspring")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

type Runner struct {
	Path string
  Arguments []string
	command *Cmd
}

func NewRunner(path string, arguments ...string) *Runner {
	return &Runner{
		Path: path,
		Arguments: arguments,
	}
}

func (runner *Runner) AddArguments (arguments ...string) {
	runner.Arguments = append(runner.Arguments, arguments...)
}

func (runner *Runner) Command(additonalArguments ...string) *Cmd {
	arguments := append(runner.Arguments, additonalArguments...)
	runner.command = Command(runner.Path, arguments...)

	return runner.command
}

func (runner *Runner) PipeCommand(additonalArguments ...string) (*Cmd, io.WriteCloser) {
	runner.Command(additonalArguments...)
	stdin := runner.StdinPipe()

	return runner.command, stdin
}


func (runner *Runner) StdinPipe() io.WriteCloser {
	stdin, err := runner.command.StdinPipe()

	if err != nil {
		log.Panic(err)
	}

	return stdin
}



