package command

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

type Meta struct {
	UI        cli.Ui
	testStdin io.Reader
}

func (m *Meta) FlagSet(name string) *pflag.FlagSet {
	f := pflag.NewFlagSet(name, pflag.ContinueOnError)

	// Create an io.Writer that writes to our UI properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.UI.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

func (m *Meta) ReadFile(path string) ([]byte, bool) {

	var rawBytes []byte
	var err error

	if path == "-" {

		var reader io.Reader = os.Stdin

		if m.testStdin != nil {
			reader = m.testStdin
		}

		rawBytes, err = ioutil.ReadAll(reader)

		if err != nil {
			m.UI.Error(fmt.Sprintf("Failed to read stdin: %v", err))
			return nil, false
		}

	} else {
		rawBytes, err = ioutil.ReadFile(path)
		if err != nil {
			m.UI.Error(fmt.Sprintf("Failed to read file: %v", err))
			return nil, false
		}
	}

	return rawBytes, true
}
