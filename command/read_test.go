package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

var testXml string = `
<books>
	<book id="1" type="short" />
	<book id="2" type="short" />
	<book id="3" type="long" />
	<book id="4" type="long" />
</books>`

func TestReadCommand_Implements(t *testing.T) {
	t.Parallel()
	var _ cli.Command = &ReadCommand{}
}

func TestReadCommand_Fails(t *testing.T) {
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

func TestReadCommand_Run(t *testing.T) {
	t.Parallel()

	ui := new(cli.MockUi)
	cmd := &ReadCommand{Meta: Meta{UI: ui}}

	file, err := ioutil.TempFile("", "xq-source")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(testXml)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if code := cmd.Run([]string{"//book/@id", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "1\n2\n3\n4\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"count(//book))", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "4\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"count(//book[@type=\"short\"]))", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "2\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()
}

func TestReadCommand_Run_stdin(t *testing.T) {
	t.Parallel()

	stdinR, stdinW, _ := os.Pipe()

	go func() {
		stdinW.WriteString(testXml)
		stdinW.Close()
	}()

	ui := new(cli.MockUi)
	cmd := &ReadCommand{Meta: Meta{UI: ui, testStdin: stdinR}}

	if code := cmd.Run([]string{"//book/@id", "-"}); code != 0 {
		assert.Equal(t, 0, code)
	}

	assert.Equal(t, "1\n2\n3\n4\n", ui.OutputWriter.String())
}
