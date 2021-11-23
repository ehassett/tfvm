package command

import "github.com/mitchellh/cli"

// Meta is a struct that contains necessary metadata used by commands.
type Meta struct {
	TerraformVersion string
	InstallPath      string
	BinPath          string
	TempPath         string
	Extension        string
	Ui 							 cli.Ui
}
