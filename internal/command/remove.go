package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tfvm/internal/helper"
)

// RemoveCommand is a Command that removes a specified installed version of Terraform.
type RemoveCommand struct {
	Meta
}

func (c *RemoveCommand) Run(args []string) int {
	if len(args) < 1 {
		err := errors.New("invalid terraform version, run `tfvm list` for a list of installed versions")
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return 1
	}

	err := removeVersion(c.TerraformVersion, c.InstallPath, c.BinPath, c.Extension, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return 1
	}

	return 0
}

func (c *RemoveCommand) Synopsis() string {
	return "Remove a specific version of Terraform"
}

func (c *RemoveCommand) Help() string {
	helpText := `
Usage: tfvm remove <version>

	Removes a specific Terraform version from the system.

	For a list of installed versions, run:
		terraform list
	`

	return strings.TrimSpace(helpText)
}

func removeVersion(
	currentVersion string,
	installPath string,
	binPath string,
	extension string,
	version string,
) error {
	// Check if version is installed.
	err := helper.IsInstalledVersion(installPath, extension, version)
	if err != nil {
		return err
	}

	// Remove the version from the install path.
	err = os.Remove(installPath + string(filepath.Separator) + "terraform" + version + extension)
	if err != nil {
		return err
	}

	// Remove the binary from the binary path if it is the current version.
	if version == currentVersion {
		err := os.Remove(binPath + string(filepath.Separator) + "terraform")
		if err != nil {
			return err
		}
	}

	fmt.Printf("Terraform v%s was successfully removed.", version)
	return nil
}
