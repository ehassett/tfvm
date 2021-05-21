// TODO: Proper error handling
// TODO: Clean up Install function
// TODO: Verify version arg in commands
// TODO: GetLatestVersion fucntion

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var basePath = os.Getenv("HOME") + "/.tfvm"
var installPath = basePath + "/versions"
var binPath = basePath + "/bin"
var currentVersion string

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

	// Find current terraform version
	if _, err := os.Stat(binPath + "/terraform"); os.IsNotExist(err) {
		currentVersion = ""
	} else {
		out, err := exec.Command(binPath+"/terraform", "-v").Output()
		if err != nil {
			panic(err)
		}
		tmp := strings.Split(string(out), "v")[1]
		currentVersion = strings.Split(tmp, "\n")[0]
	}
}

func main() {
	// Verify that a Subcommand has been provided
	if len(os.Args) < 2 {
		Help()
		os.Exit(0)
	}

	// Parse input
	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			Install(GetLatestVersion())
		} else if os.Args[2] == "list" {
			versions := GetAvailableVersions()
			for i := 0; i < len(versions); i++ {
				fmt.Println(versions[i])
			}
		} else {
			if IsAvailableVersion(os.Args[2]) {
				Install(os.Args[2])
			} else {
				fmt.Printf("Please enter a valid terraform version. For a list of available versions, run `tfvm install list`.")
			}
		}
	case "select":
		if len(os.Args) < 3 {
			fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
			os.Exit(0)
		} else {
			if IsInstalledVersion(os.Args[2]) {
				Select(os.Args[2])
			} else {
				fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
				os.Exit(0)
			}
		}
	case "list":
		versions := GetInstalledVersions()
		for i := 0; i < len(versions); i++ {
			if versions[i] == currentVersion {
				fmt.Println("* " + versions[i])
			} else {
				fmt.Println("  " + versions[i])
			}
		}
		os.Exit(0)
	case "remove":
		if len(os.Args) < 3 {
			fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
			os.Exit(0)
		} else {
			if IsInstalledVersion(os.Args[2]) {
				Remove(os.Args[2])
				os.Exit(0)
			} else {
				fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
				os.Exit(0)
			}
		}
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
	err = DownloadZip(url)
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
	Select(version)
	os.Exit(0)
}

func Select(version string) {
	if version == currentVersion {
		os.Exit(0)
	} else {
		_, err := os.Stat(binPath + "/terraform")
		if os.IsNotExist(err) {
			err = os.Symlink(installPath+"/terraform"+version, binPath+"/terraform")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Now using terraform v%s!", version)
		} else {
			err = os.Remove(binPath + "/terraform")
			if err != nil {
				panic(err)
			} else {
				err = os.Symlink(installPath+"/terraform"+version, binPath+"/terraform")
				if err != nil {
					panic(err)
				}
				fmt.Printf("Now using terraform v%s!", version)
			}
		}
	}

	os.Exit(0)
}

func GetInstalledVersions() []string {
	var versions []string

	files, err := ioutil.ReadDir(installPath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		version := strings.Trim(f.Name(), "terraform")
		versions = append(versions, version)
	}

	return versions
}

func Remove(version string) {
	if version == currentVersion {
		err := os.Remove(binPath + "/terraform")
		if err != nil {
			panic(err)
		}
		err = os.Remove(installPath + "/terraform" + version)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Terraform v%s was successfully removed.", version)
	} else {
		err := os.Remove(installPath + "/terraform" + version)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Terraform v%s was successfully removed.", version)
	}
}

func DownloadZip(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create("/tmp/tfvm.zip")
	if err != nil {
		return err
	}
	defer out.Close()

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

func GetLatestVersion() string {
	return GetAvailableVersions()[0]
}

func GetAvailableVersions() []string {
	var versions []string
	url := "https://releases.hashicorp.com/terraform/"

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		if strings.Contains(item.Text(), "terraform") {
			version := strings.Split(item.Text(), "_")[1]
			versions = append(versions, version)
		}
	})

	return versions
}

func IsAvailableVersion(version string) bool {
	versions := GetAvailableVersions()
	for _, v := range versions {
		if v == version {
			return true
		}
	}
	return false
}

func IsInstalledVersion(version string) bool {
	versions := GetInstalledVersions()
	for _, v := range versions {
		if v == version {
			return true
		}
	}
	return false
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
