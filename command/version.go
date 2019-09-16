package command

import "xq/version"

type VersionCommand struct {
	Meta
	Version *version.VersionInfo
}

func (c *VersionCommand) Help() string {
	return "Prints the xq version"
}

func (c *VersionCommand) Name() string { return "version" }

func (c *VersionCommand) Run(_ []string) int {
	c.UI.Output(c.Version.FullVersionNumber(true))
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the xq version"
}
