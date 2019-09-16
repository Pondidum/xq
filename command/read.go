package command

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/posener/complete"
)

type ReadCommand struct {
	Meta
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: xq read [options] <xpath> <file_path>

  Runs the specified <xpath> query agains the specified <file_path>.

  If the supplied path is "-", then the file is read from stdin.  Otherwise
  it is read from the path specified.

Read Options:
`

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Synopsis() string {
	return "Queries an XML document from a file or stdin"
}

func (c *ReadCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *ReadCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReadCommand) Name() string { return "read" }

func (c *ReadCommand) Run(args []string) int {

	flags := c.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output(c.Help()) }

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

	var rawBytes []byte

	if file == "-" {
		rawBytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed to read stdin: %v", err))
			return 1
		}
	} else {
		rawBytes, err = ioutil.ReadFile(file)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed to read file: %v", err))
			return 1
		}
	}

	doc, err := xmlquery.Parse(bytes.NewReader(rawBytes))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse XML input: %v", err))
		return 1
	}

	list := xmlquery.Find(doc, exp.String())

	for _, node := range list {
		c.UI.Output(node.InnerText())
	}

	return 0
}