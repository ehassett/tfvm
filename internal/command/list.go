package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethanhassett/tfvm/internal/helper"
)

// ListCommand is a Command that lists all installed versions of Terraform.
type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	versions, err := helper.GetInstalledVersions(c.InstallPath, c.Extension)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	for i := 0; i < len(versions); i++ {
		if versions[i] == c.TerraformVersion {
			fmt.Println("* " + versions[i])
		} else {
			fmt.Println("  " + versions[i])
		}
	}
	return 0
}

func (c *ListCommand) Synopsis() string {
	return "List all installed versions of Terraform"
}

func (c *ListCommand) Help() string {
	helpText := `
Usage: tfvm list

	Lists all installed Terraform versions.
	The currently selected version will be indicated with *.
	`

	return strings.TrimSpace(helpText)
}
