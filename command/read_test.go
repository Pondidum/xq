package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestReadCommand_Implements(t *testing.T) {
	t.Parallel()
	var _ cli.Command = &ReadCommand{}
}

func TestReadCommand_Run(t *testing.T) {
	t.Parallel()

	ui := new(cli.MockUi)
	cmd := &ReadCommand{Meta: Meta{UI: ui}}

	//fails on missues
	if code := cmd.Run([]string{}); code != 1 {
		assert.Equal(t, 1, code)
	}
	assert.Empty(t, ui.OutputWriter.String())
	assert.Contains(t, ui.ErrorWriter.String(), "This command takes exactly two arguments: <xpath> <file_path>\n")
	ui.ErrorWriter.Reset()

	if code := cmd.Run([]string{"//testsuite", "/this/doesnt/exist"}); code != 1 {
		assert.Equal(t, 1, code)
	}
	assert.Empty(t, ui.OutputWriter.String())
	assert.Contains(t, ui.ErrorWriter.String(), "Failed to read file: open /this/doesnt/exist: no such file or directory\n")
	ui.ErrorWriter.Reset()

	file, err := ioutil.TempFile("", "xq-source")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(file.Name())

	//fails on a bad xpath
	if code := cmd.Run([]string{"?)((*&)", file.Name()}); code != 1 {
		assert.Equal(t, 1, code)
	}
	assert.Empty(t, ui.OutputWriter.String())
	assert.Contains(t, ui.ErrorWriter.String(), "Failed to parse xpath: ?)((*&) has an invalid token.\n")
	ui.ErrorWriter.Reset()

}
