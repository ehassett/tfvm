package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
)

// TestRemove sets up the filesystem and Meta and tests various RemoveCommand cases.
func TestRemove(t *testing.T) {
	workDir, err := ioutil.TempDir("", "tfvm-test-command-remove")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(workDir)

	installDir, err := ioutil.TempDir(workDir, "versions")
	if err != nil {
		t.Fatalf("cannot create versions directory: %s", err)
	}

	binDir, err := ioutil.TempDir(workDir, "bin")
	if err != nil {
		t.Fatalf("cannot create bin directory: %s", err)
	}

	binFileName := binDir + string(filepath.Separator) + "terraform"
	binFile, err := os.Create(binFileName)
	if err != nil {
		t.Fatalf("cannot create stub binary: %s", err)
	}

	currentVerFileName := installDir + string(filepath.Separator) + "terraform1.0.0"
	currentVerFile, err := os.Create(currentVerFileName)
	if err != nil {
		t.Fatalf("cannot create stub version file: %s", err)
	}

	verFileName := installDir + string(filepath.Separator) + "terraform0.15.0"
	verFile, err := os.Create(verFileName)
	if err != nil {
		t.Fatalf("cannot create stub version file: %s", err)
	}

	removeTestCase := func(test func(t *testing.T, c *RemoveCommand, ui *cli.MockUi)) func(t *testing.T) {
		return func(t *testing.T) {
			ui := new(cli.MockUi)

			c := &RemoveCommand{
				Meta: Meta{
					TerraformVersion: "1.0.0",
					InstallPath:      installDir,
					BinPath:          binDir,
					TempPath:         "tfvm.zip",
					Extension:        "",
					Ui:               ui,
				},
			}

			test(t, c, ui)
		}
	}

	// Pass in the current installed version and expect removal of both the bin/ and versions/ files.
	t.Run("current terraform version", removeTestCase(func(t *testing.T, c *RemoveCommand, ui *cli.MockUi) {
		status := c.Run([]string{"1.0.0"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); !os.IsNotExist(err) {
			t.Fatalf("failed to remove current executable\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); !os.IsNotExist(err) {
			t.Fatalf("failed to remove current version file\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in an installed version and expect removal of the versions/ file.
	t.Run("installed terraform version", removeTestCase(func(t *testing.T, c *RemoveCommand, ui *cli.MockUi) {
		status := c.Run([]string{"0.15.0"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(verFile.Name()); !os.IsNotExist(err) {
			t.Fatalf("failed to remove current version file\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in a not installed version and expect an error.
	t.Run("not installed terraform version", removeTestCase(func(t *testing.T, c *RemoveCommand, ui *cli.MockUi) {
		status := c.Run([]string{"0.13.5"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))

	// Pass in an invalid version and expect an error.
	t.Run("not valid terraform version", removeTestCase(func(t *testing.T, c *RemoveCommand, ui *cli.MockUi) {
		status := c.Run([]string{"invalid_version"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))

	// Pass in no version and expect an error.
	t.Run("not valid terraform version", removeTestCase(func(t *testing.T, c *RemoveCommand, ui *cli.MockUi) {
		status := c.Run([]string{""})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))
}
