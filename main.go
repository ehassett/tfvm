package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var currentVersion, basePath, installPath, binPath, zipPath, extension string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	basePath = home + string(filepath.Separator) + ".tfvm"
	installPath = basePath + string(filepath.Separator) + "versions"
	binPath = basePath + string(filepath.Separator) + "bin"
	zipPath = basePath + string(filepath.Separator) + "tfvm.zip"

	switch runtime.GOOS {
	case "windows":
		extension = ".exe"
	case "linux":
		extension = ""
	default:
		extension = ""
		err := errors.New("operating system could not be verified")
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.Mkdir(basePath, 0755)
	}
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.Mkdir(installPath, 0755)
	}
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	if _, err := os.Stat(binPath + string(filepath.Separator) + "terraform" + extension); os.IsNotExist(err) {
		currentVersion = ""
	} else {
		out, err := exec.Command(binPath+string(filepath.Separator)+"terraform"+extension, "-v").Output()
		if err != nil {
			panic(err)
		}
		tmp := strings.Split(string(out), "v")[1]
		currentVersion = strings.Split(tmp, "\n")[0]
	}
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			latest, err := getLatestVersion()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			install(latest)
		}

		if os.Args[2] == "list" {
			versions, err := getAvailableVersions()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			for i := 0; i < len(versions); i++ {
				fmt.Println(versions[i])
			}
			os.Exit(0)
		}
		install(os.Args[2])
	case "select":
		if len(os.Args) < 3 {
			fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
			os.Exit(0)
		}
		selectVersion(os.Args[2])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
			os.Exit(0)
		}
		remove(os.Args[2])
	case "list":
		list()
	default:
		help()
	}
}
