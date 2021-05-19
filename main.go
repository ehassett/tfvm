// TODO: Proper error handling
// TODO: Clean up Install function
// TODO: Verify version arg in commands
// TODO: GetLatestVersion fucntion

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var basePath = os.Getenv("HOME") + "/.tfvm"
var installPath = basePath + "/versions"
var binPath = basePath + "/bin"

func init() {
	// Create file structure if needed
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, 0755)
	}

	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.Mkdir(installPath, 0755)
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}
}

func main() {
	// Verify that a Subcommand has been provided
	if len(os.Args) < 2 {
		Help()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			Install(GetLatestVersion())
		} else {
			Install(os.Args[2])
		}
	case "select":
		if len(os.Args) < 3 {
			fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run tfvm list.")
		}
		SetTerraformVersion(os.Args[2])
	default:
		Help()
	}
}

func Help() {
	fmt.Printf("tfvm usage:\n")
	fmt.Printf("\thelp\n")
	fmt.Printf("\t\tShow this help text.\n")
	fmt.Printf("\tinstall\n")
	fmt.Printf("\t\tUse `tfvm install [version]` to install terraform. If no version is specified, the latest will be installed.\n")
}

func Install(version string) {
	// Check if version is already installed
	_, err := os.Stat(installPath + "/terraform" + version)
	if !os.IsNotExist(err) {
		fmt.Printf("Terraform v%s is already installed. Run tfvm select %s to use a new version.", version, version)
		os.Exit(0)
	}

	url := "https://releases.hashicorp.com/terraform/" + version + "/terraform_" + version + "_" + GetArchitecture() + ".zip"

	fmt.Printf("Installing terraform v%s...\n", version)
	// Download terraform zip
	fmt.Printf("Downloading terraform v%s from %s\n", version, url)
	err = DownloadTerraform(url)
	if err != nil {
		panic(err)
	}
	// Extract zip to install path
	err = Unzip("/tmp/tfvm.zip", installPath)
	if err != nil {
		panic(err)
	}
	// Remove downloaded zip
	err = os.Remove("/tmp/tfvm.zip")
	if err != nil {
		panic(err)
	}
	// Add version to file name
	err = os.Rename(installPath+"/terraform", installPath+"/terraform"+version)
	if err != nil {
		panic(err)
	}
	// Set selected version to the newly installed bin
	err = SetTerraformVersion(version)
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}

func DownloadTerraform(url string) error {
	// Get data
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create file
	out, err := os.Create("/tmp/tfvm.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write body to file
	_, err = io.Copy(out, response.Body)
	return err
}

func Unzip(src string, dest string) error {

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

func SetTerraformVersion(version string) error {
	_, err := os.Stat(binPath + "/terraform")
	if os.IsNotExist(err) {
		err = os.Symlink(installPath+"/terraform"+version, binPath+"/terraform")
		if err != nil {
			return err
		}
		fmt.Printf("Now using terraform v%s!", version)
	} else {
		err = os.Remove(binPath + "/terraform")
		if err != nil {
			return err
		} else {
			err = os.Symlink(installPath+"/terraform"+version, binPath+"/terraform")
			if err != nil {
				return err
			}
			fmt.Printf("Now using terraform v%s!", version)
		}
	}

	return nil
}

func GetLatestVersion() string {
	return "0.13.5"
}

func GetArchitecture() string {
	var arch string

	if runtime.GOOS == "linux" {
		switch runtime.GOARCH {
		case "386":
			arch = "linux_386"
		case "amd64":
			arch = "linux_amd64"
		case "arm":
			arch = "linux_arm"
		case "arm64":
			arch = "linux_arm64"
		}
	} else {
		fmt.Printf("Could not verify your OS or architecture. Aborting.")
		os.Exit(1)
	}

	return arch
}
