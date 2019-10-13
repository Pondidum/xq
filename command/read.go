package command

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

type ReadCommand struct {
	Meta
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: xq read [options] <xpath> <file_path>

  Runs the specified <xpath> query against the specified <file_path>.

  If the supplied path is "-", then the file is read from stdin.  Otherwise
  it is read from the path specified.

Read Options:

  --output
    Optional. Specifies the format to output the xpath result in.
    Valid values: xml, raw. Defaults to xml.
`

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Synopsis() string {
	return "Queries an XML document from a file or stdin"
}

func (c *ReadCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"--output": complete.PredictSet("raw", "xml"),
	}
}

func (c *ReadCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReadCommand) Name() string { return "read" }

func (c *ReadCommand) Run(args []string) int {

	var output string

	flags := c.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output(c.Help()) }
	flags.StringVar(&output, "output", "xml", "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if l := len(args); l != 2 {
		c.UI.Error("This command takes exactly two arguments: <xpath> <file_path>")
		return 1
	}

	query := args[0]
	file := args[1]

	exp, err := xpath.Compile(query)

	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse xpath: %v", err))
		return 1
	}

	rawBytes, ok := c.ReadFile(file)
	if ok == false {
		return 1
	}

	doc, err := xmlquery.Parse(bytes.NewReader(rawBytes))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse XML input: %v", err))
		return 1
	}

	result := exp.Evaluate(xmlquery.CreateXPathNavigator(doc))
	iterator, ok := result.(*xpath.NodeIterator)

	if ok {
		for iterator.MoveNext() {
			navigator := iterator.Current()
			writeOutput(c.UI, output, navigator)
		}
	} else {
		c.UI.Output(fmt.Sprintf("%+v", result))
	}

	return 0
}

func writeOutput(ui cli.Ui, mode string, navigator xpath.NodeNavigator) {

	switch mode {
	case "xml":
		var buffer bytes.Buffer
		render(&buffer, navigator)

		ui.Output(buffer.String())
	case "raw":
		ui.Output(navigator.Value())
	}

}

func render(buffer *bytes.Buffer, nav xpath.NodeNavigator) {

	if nav.NodeType() == xpath.TextNode {
		xml.EscapeText(buffer, []byte(strings.TrimSpace(nav.Value())))
		return
	}

	name := nav.LocalName()

	if nav.NodeType() == xpath.AttributeNode {
		buffer.WriteString(fmt.Sprintf("<%s>%s</%s>", name, nav.Value(), name))
		return
	}

	buffer.WriteString("<")
	buffer.WriteString(name)

	withAttributes := false
	for nav.MoveToNextAttribute() {
		buffer.WriteString(fmt.Sprintf(` %s="%s"`, nav.LocalName(), nav.Value()))
		withAttributes = true
	}

	if withAttributes {
		nav.MoveToParent()
	}

	if nav.MoveToChild() == false {
		buffer.WriteString(" />")
		return
	}

	buffer.WriteString(">")

	render(buffer, nav)

	for nav.MoveToNext() {
		render(buffer, nav)
	}

	buffer.WriteString("</" + name + ">")
	nav.MoveToParent()
}
