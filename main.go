package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/cli"
)

var appVersion string = "2.0.0"

type ui struct {
	cli.Ui
}

func (ui *ui) Warn(msg string) {
	ui.Ui.Output(msg)
}

var Ui = &ui{&cli.BasicUi{
	Writer:      os.Stdout,
	ErrorWriter: os.Stderr,
	Reader:      os.Stdin,
}}

func main() {
	c := cli.NewCLI("tfvm", appVersion)
	c.Args = os.Args[1:]
	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		Ui.Error(fmt.Sprintf("Error executing CLI: %s", err.Error()))
	}
	os.Exit(exitStatus)
}

func init() {
	var terraformVersion, basePath, installPath, binPath, tempPath, extension string

	// Determine paths and extensions based on OS.
	home, err := os.UserHomeDir()
	if err != nil {
		Ui.Error(fmt.Sprintf("Failed to find home directory: %s", err))
		os.Exit(1)
	}
	basePath = home + string(filepath.Separator) + ".tfvm"
	installPath = basePath + string(filepath.Separator) + "versions"
	binPath = basePath + string(filepath.Separator) + "bin"
	tempPath = basePath + string(filepath.Separator) + "tfvm.zip"

	switch runtime.GOOS {
	case "windows":
		extension = ".exe"
	case "linux":
		extension = ""
	case "darwin":
		extension = ""
	default:
		extension = ""
		err := errors.New("operating system could not be verified")
		Ui.Error(fmt.Sprintf("Failed to determine extension: %s", err))
		os.Exit(1)
	}

	// Create directory structure if needed.
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, 0755)
	}
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.Mkdir(installPath, 0755)
	}
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	// Set current Terraform version if set.
	if _, err := os.Stat(binPath + string(filepath.Separator) + "terraform" + extension); os.IsNotExist(err) {
		terraformVersion = ""
	} else {
		out, err := exec.Command(binPath+string(filepath.Separator)+"terraform"+extension, "-v").Output()
		if err != nil {
			Ui.Error(fmt.Sprintf("Failed to determine current terraform version: %s", err))
			os.Exit(1)
		}
		tmp := strings.Split(string(out), "v")[1]
		terraformVersion = strings.Split(tmp, "\n")[0]
	}

	// Pass initialized values to initCommands for Meta.
	initCommands(terraformVersion, installPath, binPath, tempPath, extension, Ui)
}
