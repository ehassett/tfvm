package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
)

// TestInstall sets up the filesystem and Meta and tests various InstallCommand cases.
func TestInstall(t *testing.T) {
	workDir, err := ioutil.TempDir("", "tfvm-test-command-install")
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

	installTestCase := func(test func(t *testing.T, c *InstallCommand, ui *cli.MockUi)) func(t *testing.T) {
		return func(t *testing.T) {
			ui := new(cli.MockUi)

			c := &InstallCommand{
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

	// Pass in the current version and expect graceful error.
	t.Run("current terraform version", installTestCase(func(t *testing.T, c *InstallCommand, ui *cli.MockUi) {
		status := c.Run([]string{"1.0.0"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the executable\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in a valid version and expect it to be installed.
	t.Run("valid terraform version", installTestCase(func(t *testing.T, c *InstallCommand, ui *cli.MockUi) {
		status := c.Run([]string{"1.0.2"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(installDir + string(filepath.Separator) + "terraform1.0.2"); os.IsNotExist(err) {
			t.Fatalf("failed to install new version\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in an invalid version and expect error.
	t.Run("invalid terraform version", installTestCase(func(t *testing.T, c *InstallCommand, ui *cli.MockUi) {
		status := c.Run([]string{"9999.9999.9999"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the executable\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in an no version and expect the latest to be installed.
	t.Run("no specified version", installTestCase(func(t *testing.T, c *InstallCommand, ui *cli.MockUi) {
		status := c.Run([]string{})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the executable\nstderr: %s", ui.ErrorWriter.String())
		}
	}))
}
