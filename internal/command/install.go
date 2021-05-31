package command

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"tfvm/internal/helper"
)

// InstallCommand is a Command that downloads and installs a specific Terraform version.
type InstallCommand struct {
	Meta
}

func (c *InstallCommand) Run(args []string) int {
	if len(args) < 1 {
		err := installLatest(c.TerraformVersion, c.InstallPath, c.BinPath, c.TempPath, c.Extension)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}
		return 0
	}

	switch args[0] {
	case "--list", "-list", "-l":
		versions, err := helper.GetAvailableVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}
		for i := 0; i < len(versions); i++ {
			fmt.Println(versions[i])
		}
	default:
		err := installVersion(c.TerraformVersion, c.InstallPath, c.BinPath, c.TempPath, c.Extension, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}
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

	For a list of available versions, run:
  	tfvm install --list

	Options:
		--list, -l	List available versions of Terraform
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
	// Check if the selected version is already installed.
	_, err := os.Stat(installPath + string(filepath.Separator) + "terraform" + version + extension)
	if !os.IsNotExist(err) {
		err = errors.New("version already installed, run `tfvm use " + version + "` to use this version")
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
	fmt.Printf("Downloading terraform v%s from %s\n", version, url)
	err = downloadArchive(url, tempPath)
	if err != nil {
		return err
	}

	fmt.Println("Extracting archive...")
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

	fmt.Printf("Terraform v%s successfully installed. Run `tfvm use %s` to use this new version.", version, version)
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
			return fmt.Errorf("%s: illegal file path", fpath)
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
