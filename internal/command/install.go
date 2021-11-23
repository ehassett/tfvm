package command

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethanhassett/tfvm/internal/helper"
)

// InstallCommand is a Command that downloads and installs a specific Terraform version.
type InstallCommand struct {
	Meta
}

func (c *InstallCommand) Run(args []string) int {
	if len(args) < 1 {
		err := installLatest(c.TerraformVersion, c.InstallPath, c.BinPath, c.TempPath, c.Extension)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Could not install latest version: %s", err))
			return 1
		}
		return 0
	}

	switch args[0] {
	case "--list", "-list", "-l":
		versions, err := helper.GetAvailableVersions()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Could not show available versions: %s", err))
			return 1
		}

		pager := os.Getenv("PAGER")
		if pager != "" {
			cmd := exec.Command(pager)
			cmd.Stdin = strings.NewReader(strings.Join(versions, "\n"))
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Error using $PAGER: %s", err))
			}
		} else {
			for i := 0; i < len(versions); i++ {
				c.Ui.Output(versions[i])
			}
		}

	default:
		err := installVersion(c.TerraformVersion, c.InstallPath, c.BinPath, c.TempPath, c.Extension, args[0])
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Could not install specified version: %s", err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Terraform v%s successfully installed. Run `tfvm use %s` to use this new version.", args[0], args[0]))
	}
	return 0
}

func (c *InstallCommand) Synopsis() string {
	return "Install a version of Terraform"
}

func (c *InstallCommand) Help() string {
	helpText := `
Usage: tfvm install [version]

	Installs a Terraform binary according to the specified version.
	If no version is specified, tfvm will default to the latest available version.
	Version specification can be to the patch or minor version.
	Only specifying a minor version will install the latest patch of that version.

	For a list of available versions, run:
  	tfvm install --list

	Options:
		--list, -l	List available versions of Terraform

	Examples:
		tfvm install 1.0.0	Installs Terraform v1.0.0
		tfvm install 1.0	Installs the latest of Terraform v1.0.x
	`

	return strings.TrimSpace(helpText)
}

// installVersion verifies and installs the specified version of Terraform.
func installVersion(
	currentVersion string,
	installPath string,
	binPath string,
	tempPath string,
	extension string,
	version string,
) error {
	if strings.Count(version, ".") == 1 {
		fullVersion, err := getMinorVersion(version)
		if err != nil {
			return err
		}
		version = fullVersion
	}

	// Check if the selected version is already installed.
	_, err := os.Stat(installPath + string(filepath.Separator) + "terraform" + version + extension)
	if !os.IsNotExist(err) {
		err = errors.New("already installed, run `tfvm use " + version + "` to use this version")
		return err
	}

	// Check if the selected version is available to install.
	err = helper.IsAvailableVersion(version)
	if err != nil {
		return err
	}

	arch, err := getArchitecture()
	if err != nil {
		return err
	}

	url := "https://releases.hashicorp.com/terraform/" + version + "/terraform_" + version + "_" + arch + ".zip"
	err = downloadArchive(url, tempPath)
	if err != nil {
		return err
	}

	err = unzipArchive(tempPath, installPath)
	if err != nil {
		return err
	}
	err = os.Remove(tempPath)
	if err != nil {
		return err
	}
	err = os.Rename(installPath+string(filepath.Separator)+"terraform"+extension, installPath+string(filepath.Separator)+"terraform"+version+extension)
	if err != nil {
		return err
	}

	return nil
}

// installLatest installs the newest available version of Terraform.
func installLatest(
	currentVersion string,
	installPath string,
	binPath string,
	tempPath string,
	extension string,
) error {
	versions, err := helper.GetAvailableVersions()
	if err != nil {
		return err
	}

	err = installVersion(currentVersion, installPath, binPath, tempPath, extension, versions[0])
	return err
}

// getArchitecture determines OS and Arch information.
func getArchitecture() (string, error) {
	var arch string = ""
	var err error = nil

	if runtime.GOOS == "windows" {
		switch runtime.GOARCH {
		case "386":
			arch = "windows_386"
		case "amd64":
			arch = "windows_amd64"
		default:
			err = errors.New("architecture could not be verified for installation")
		}
	} else if runtime.GOOS == "linux" {
		switch runtime.GOARCH {
		case "386":
			arch = "linux_386"
		case "amd64":
			arch = "linux_amd64"
		case "arm":
			arch = "linux_arm"
		case "arm64":
			arch = "linux_arm64"
		default:
			err = errors.New("architecture could not be verified for installation")
		}
	} else if runtime.GOOS == "darwin" {
		switch runtime.GOARCH {
		case "amd64":
			arch = "darwin_amd64"
		case "arm64":
			arch = "darwin_arm64"
		default:
			err = errors.New("architecture could not be verified for installation")
		}
	} else {
		err = errors.New("operating system is not supported for installation")
	}

	return arch, err
}

// downloadArchive downloads the Zip at the specified URL.
func downloadArchive(url string, tempPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzipArchive unzips the Zip at src to the path at dest.
func unzipArchive(src string, dest string) error {

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			err = errors.New("illegal file path: " + fpath)
			return err
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// getMinorVersion gets the latest Terrform version from a minor version.
func getMinorVersion(version string) (string, error) {
	versions, err := helper.GetAvailableVersions()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(versions); i++ {
		tmp := strings.Split(versions[i], ".")[0] + "." + strings.Split(versions[i], ".")[1]
		if strings.Contains(tmp, version) {
			return versions[i], nil
		}
	}

	err = errors.New("invalid minor specification, run `tfvm install --list` to see all available versions")
	return "", err
}
