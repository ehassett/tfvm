package command

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
)

// TestUse sets up the filesystem and Meta and tests various UseCommand cases.
func TestUse(t *testing.T) {
	workDir, err := ioutil.TempDir("", "tfvm-test-command-use")
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

	useTestCase := func(test func(t *testing.T, c *UseCommand, ui *cli.MockUi)) func(t *testing.T) {
		return func(t *testing.T) {
			ui := new(cli.MockUi)

			c := &UseCommand{
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

	// Pass in an installed but not current version and expect no errors or file changes.
	t.Run("installed terraform version", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		status := c.Run([]string{"0.15.0"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep an executable\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(verFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep the selected version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in the current version and expect no errors or file changes.
	t.Run("current terraform version", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		status := c.Run([]string{"1.0.0"})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep an executable\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(verFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep the selected version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to keep the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in an invalid version and expect an error.
	t.Run("invalid terraform version", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		status := c.Run([]string{"invalid_version"})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))

	// Pass in no version and expect an error.
	t.Run("no terraform version", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		status := c.Run([]string{})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))

	// Pass in no version and expect to use valid .tfversion file with no file changes.
	t.Run("use valid .tfversion", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory\nstderr: %s", ui.ErrorWriter.String())
		}

		tfversionFile, err := os.Create(cwd + string(filepath.Separator) + ".tfversion")
		if err != nil {
			t.Fatalf("cannot create stub .tfversion file: %s", err)
		}
		defer os.Remove(tfversionFile.Name())

		w := bufio.NewWriter(tfversionFile)
		_, err = w.WriteString("0.15.0\n")
		if err != nil {
			t.Fatalf("cannot write stub .tfversion file: %s", err)
		}
		w.Flush()

		if _, err := os.Stat(tfversionFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to stat .tfversion in cwd\nstderr: %s", ui.ErrorWriter.String())
		}

		status := c.Run([]string{})
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}

		if _, err := os.Stat(binFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the executable\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(verFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the selected version file\nstderr: %s", ui.ErrorWriter.String())
		}

		if _, err := os.Stat(currentVerFile.Name()); os.IsNotExist(err) {
			t.Fatalf("unexpectedly removed the previous version file\nstderr: %s", ui.ErrorWriter.String())
		}
	}))

	// Pass in no version and expect to error on  imvalid .tfversion file.
	t.Run("use invalid .tfversion", useTestCase(func(t *testing.T, c *UseCommand, ui *cli.MockUi) {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory\nstderr: %s", ui.ErrorWriter.String())
		}

		tfversionFile, err := os.Create(cwd + string(filepath.Separator) + ".tfversion")
		if err != nil {
			t.Fatalf("cannot create stub .tfversion file: %s", err)
		}
		defer os.Remove(tfversionFile.Name())

		if _, err := os.Stat(tfversionFile.Name()); os.IsNotExist(err) {
			t.Fatalf("failed to stat .tfversion in cwd\nstderr: %s", ui.ErrorWriter.String())
		}

		status := c.Run([]string{})
		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr: %s", status, ui.ErrorWriter.String())
		}
	}))
}
