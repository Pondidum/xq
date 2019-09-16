package command

import (
	"bufio"
	"io"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

type Meta struct {
	UI cli.Ui
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
