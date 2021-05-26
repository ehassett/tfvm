package main

import (
	"tfvm/internal/command"

	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available tfvm commands.
var Commands map[string]cli.CommandFactory

func initCommands(
	terraformVersion string,
	installPath string,
	binPath string,
	tempPath string,
	extension string,
) {

	meta := command.Meta{
		TerraformVersion: terraformVersion,
		InstallPath:      installPath,
		BinPath:          binPath,
		TempPath:         tempPath,
		Extension:        extension,
	}

	Commands = map[string]cli.CommandFactory{
		"install": func() (cli.Command, error) {
			return &command.InstallCommand{
				Meta: meta,
			}, nil
		},
		"use": func() (cli.Command, error) {
			return &command.UseCommand{
				Meta: meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: meta,
			}, nil
		},
		"remove": func() (cli.Command, error) {
			return &command.RemoveCommand{
				Meta: meta,
			}, nil
		},
	}
}
