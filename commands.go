package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func install(version string) {
	_, err := os.Stat(installPath + string(filepath.Separator) + "terraform" + version + extension)
	if !os.IsNotExist(err) {
		fmt.Printf("Terraform v%s is already installed. Run `tfvm select %s` to use this version.", version, version)
		os.Exit(0)
	}

	valid, err := isAvailableVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	if !valid {
		fmt.Printf("Please enter a valid terraform version. For a list of available versions, run `tfvm install list`.")
		os.Exit(0)
	}

	arch, err := getArchitecture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	url := "https://releases.hashicorp.com/terraform/" + version + "/terraform_" + version + "_" + arch + ".zip"

	fmt.Printf("Downloading terraform v%s from %s\n", version, url)
	err = downloadArchive(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Extracting archive...")
	err = unzipArchive(zipPath, installPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	err = os.Remove(zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	err = os.Rename(installPath+string(filepath.Separator)+"terraform"+extension, installPath+string(filepath.Separator)+"terraform"+version+extension)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Terraform v%s successfully installed.\n", version)
	selectVersion(version)
	os.Exit(0)
}

func selectVersion(version string) {
	valid, err := isInstalledVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if !valid {
		fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
		os.Exit(0)
	}

	if version == currentVersion {
		fmt.Printf("Terraform %s already in use!", version)
		os.Exit(0)
	}
	_, err = os.Stat(binPath + string(filepath.Separator) + "terraform" + extension)
	if !os.IsNotExist(err) {
		err = os.Remove(binPath + string(filepath.Separator) + "terraform" + extension)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	err = os.Link(installPath+string(filepath.Separator)+"terraform"+version+extension, binPath+string(filepath.Separator)+"terraform"+extension)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Now using terraform v%s.", version)
	os.Exit(0)
}

func remove(version string) {
	valid, err := isInstalledVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if !valid {
		fmt.Printf("Please enter a valid terraform version. For a list of installed versions, run `tfvm list`.")
		os.Exit(0)
	}

	err = os.Remove(installPath + string(filepath.Separator) + "terraform" + version + extension)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if version == currentVersion {
		err := os.Remove(binPath + string(filepath.Separator) + "terraform")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Terraform v%s was successfully removed.", version)
	os.Exit(0)
}

func list() {
	versions, err := getInstalledVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(versions); i++ {
		if versions[i] == currentVersion {
			fmt.Println("* " + versions[i])
		} else {
			fmt.Println("  " + versions[i])
		}
	}
	os.Exit(0)
}

func help() {
	fmt.Println("tfvm usage:")
	fmt.Println("  help")
	fmt.Println("    tfvm help - Shows this help text.")
	fmt.Println("  install")
	fmt.Println("    tfvm install [version] - Installs terraform. If no version is specified, the latest will be installed.")
	fmt.Println("    tfvm install list - Lists the available terraform versions.")
	fmt.Println("  list")
	fmt.Println("    tfvm list - Lists all installed terraform versions. The current version is indicated with a *.")
	fmt.Println("  remove")
	fmt.Println("    tfvm remove <version> - Uninstalls the specified terraform version.")
	fmt.Println("  select")
	fmt.Println("    tfvm select <version> - Selects the specified terraform version to be used.")
	os.Exit(0)
}
