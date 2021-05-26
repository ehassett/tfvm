package command

import (
	"errors"
	"fmt"
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
	if len(args) < 1 {
		err := errors.New("invalid terraform version, run `tfvm list` for a list of installed versions")
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return 1
	}
	err := useVersion(c.TerraformVersion, c.InstallPath, c.BinPath, c.Extension, args[0])
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
Usage: tfvm use <version>

	Selects a Terraform version to use.

	For a list of installed versions, run:
		terraform list
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
