package command

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
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
	assert.Contains(t, ui.ErrorWriter.String(), "Failed to read file: open /this/doesnt/exist: ")
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

func TestReadCommand_Run_stdin(t *testing.T) {
	t.Parallel()

	stdinR, stdinW, _ := os.Pipe()

	go func() {
		stdinW.WriteString(testXml)
		stdinW.Close()
	}()

	ui := new(cli.MockUi)
	cmd := &ReadCommand{Meta: Meta{UI: ui, testStdin: stdinR}}

	if code := cmd.Run([]string{"--output", "raw", "//book/@id", "-"}); code != 0 {
		assert.Equal(t, 0, code)
	}

	assert.Equal(t, "1\n2\n3\n4\n", ui.OutputWriter.String())
}

func TestReadCommand_Run_RawOutput(t *testing.T) {
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

	if code := cmd.Run([]string{"--output", "raw", "//book/@id", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "1\n2\n3\n4\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"--output", "raw", "count(//book)", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "4\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"--output", "raw", "count(//book[@type=\"short\"])", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "2\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()
}

func TestReadCommand_Run_XmlOutput(t *testing.T) {
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

	if code := cmd.Run([]string{"--output", "xml", "//book[1]", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "<book id=\"1\" type=\"short\" />\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"--output", "xml", "//book/@id", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Contains(t, "<id>1</id>\n<id>2</id>\n<id>3</id>\n<id>4</id>\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()

	if code := cmd.Run([]string{"--output", "xml", "count(//book)", file.Name()}); code != 0 {
		assert.Equal(t, 0, code)
	}
	assert.Equal(t, "4\n", ui.OutputWriter.String())
	ui.OutputWriter.Reset()
}

func renderXml(xml string) string {
	doc, _ := xmlquery.Parse(strings.NewReader(xml))
	exp, _ := xpath.Compile("/*")
	result := exp.Evaluate(xmlquery.CreateXPathNavigator(doc))
	iterator, _ := result.(*xpath.NodeIterator)

	iterator.MoveNext()
	nav := iterator.Current()

	var buf bytes.Buffer
	render(&buf, nav)

	return buf.String()
}

func TestRender(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `<testing />`, renderXml(`<testing />`))
	assert.Equal(t, `<testing id="one" />`, renderXml(`<testing id="one" />`))
	assert.Equal(t, `<testing id="one" name="two" />`, renderXml(`<testing id="one" name="two" />`))
	assert.Equal(t, `<testing />`, renderXml(`<testing></testing>`))

	assert.Equal(t, `<testing>plain text</testing>`, renderXml(`<testing>plain text</testing>`))
	assert.Equal(t, `<testing><child /></testing>`, renderXml(`<testing><child /></testing>`))
	assert.Equal(t, `<testing id="parent"><child /></testing>`, renderXml(`<testing id="parent"><child /></testing>`))
	assert.Equal(t, `<testing><child id="1" /><child id="2" /></testing>`, renderXml(`<testing><child id="1" /><child id="2" /></testing>`))
	assert.Equal(t, `<testing name="test"><child id="1" /><child id="2" /></testing>`, renderXml(`<testing name="test"><child id="1" /><child id="2" /></testing>`))

	assert.Equal(t, `<books><book><name lang="en">first</name><title lang="en">the title</title></book></books>`, renderXml(`<books><book><name lang="en">first</name><title lang="en">the title</title></book></books>`))
	assert.Equal(t, `<books><book><name lang="en">first</name><title lang="en">the title</title></book><book><name lang="en">second</name><title lang="en">different</title></book></books>`, renderXml(`<books><book><name lang="en">first</name><title lang="en">the title</title></book><book><name lang="en">second</name><title lang="en">different</title></book></books>`))
}
