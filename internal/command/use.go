package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"tfvm/internal/helper"
)

// UseCommand is a Command that selects a Terraform version to be used.
type UseCommand struct {
	Meta
}

func (c *UseCommand) Run(args []string) int {
	var version string

	if len(args) < 1 {
		// Get working directory.
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}

		// Test if .tfversion exists in working directory.
		if _, err := os.Stat(cwd + string(filepath.Separator) + ".tfversion"); os.IsNotExist(err) {
			err := errors.New("invalid terraform version, run `tfvm list` for a list of installed versions")
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}

		// Read .tfversion.
		version, err = getDirVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}

		fmt.Printf("Switching to v%s from .tfversion file.\n", version)
	} else {
		version = args[0]
	}

	err := useVersion(c.TerraformVersion, c.InstallPath, c.BinPath, c.Extension, version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return 1
	}
	return 0
}

func (c *UseCommand) Synopsis() string {
	return "Select a version of Terraform to use"
}

func (c *UseCommand) Help() string {
	helpText := `
Usage: tfvm use [version]

	Selects a Terraform version to use.
	If no version is specified, tfvm will try to select the version specified in .tfversion if it exists in the current directory.

	For a list of installed versions, run:
		tfvm list
	`

	return strings.TrimSpace(helpText)
}

// useVersion copies the appropriate binary version to the binPath to be used.
func useVersion(
	currentVersion string,
	installPath string,
	binPath string,
	extension string,
	version string) error {
	// Check if specified version is installed.
	err := helper.IsInstalledVersion(installPath, extension, version)
	if err != nil {
		return err
	}

	// Return if desired version is already current.
	if version == currentVersion {
		fmt.Printf("Now using terraform v%s.", version)
		return nil
	}

	// Remove binary from binary path if it exists.
	_, err = os.Stat(binPath + string(filepath.Separator) + "terraform" + extension)
	if !os.IsNotExist(err) {
		err = os.Remove(binPath + string(filepath.Separator) + "terraform" + extension)
		if err != nil {
			return err
		}
	}

	// Link new binary to binary path.
	err = os.Link(installPath+string(filepath.Separator)+"terraform"+version+extension, binPath+string(filepath.Separator)+"terraform"+extension)
	if err != nil {
		return err
	}

	fmt.Printf("Now using terraform v%s.", version)
	return nil
}

// getDirVersion reads the version from a .tfversion file.
func getDirVersion() (string, error) {
	var dirVersion string = ""

	// Open file for reading.
	f, err := os.OpenFile(".tfversion", os.O_RDWR, 0600)
	if err != nil {
		return dirVersion, err
	}
	defer f.Close()

	// Read data from file.
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return dirVersion, err
	}

	// Parse data by line.
	lines := strings.Split(string(raw), "\n")
	if len(lines) < 1 {
		err = errors.New("invalid .tfversion file")
		return dirVersion, err
	}
	dirVersion = lines[0]

	return dirVersion, nil
}
