package helper

import (
	"errors"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetAvailableVersions returns a list of currently available Terraform versions.
func GetAvailableVersions() ([]string, error) {
	var versions []string
	var err error = nil
	url := "https://releases.hashicorp.com/terraform/"

	resp, err := http.Get(url)
	if err != nil {
		return versions, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return versions, err
	}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		if strings.Contains(item.Text(), "terraform") {
			version := strings.Split(item.Text(), "_")[1]

			// Filter out unsupported versions on Apple Silicon
			if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
				versionSlice := strings.Split(version, ".")
				majInt, _ := strconv.Atoi(versionSlice[0])
				patchInt, _ := strconv.Atoi(versionSlice[2])

				// Do not include pre-release versions or incompatible versions
				if majInt >= 1 && patchInt >= 2 && !strings.Contains(version, "-") {
					versions = append(versions, version)
				}
			} else {
				// Do not include pre-release versions
				if !strings.Contains(version, "-") {
					versions = append(versions, version)
				}
			}
		}
	})
	return versions, err
}

// IsAvailableVersion returns true if the specified version of Terraform is in the list of available versions.
func IsAvailableVersion(version string) error {
	var err error = nil

	// Check that Apple Silicon users select a valid version (v1.0.2+).
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		versionSlice := strings.Split(version, ".")

		majInt, err := strconv.Atoi(versionSlice[0])
		if err != nil {
			return err
		}

		patchInt, err := strconv.Atoi(versionSlice[2])
		if err != nil {
			return err
		}

		if !(majInt >= 1 && patchInt >= 2) {
			err = errors.New("only Terraform v1.0.2+ is supported on Apple Silicon")
			return err
		}
	}

	versions, err := GetAvailableVersions()
	if err != nil {
		return err
	}

	for _, v := range versions {
		if v == version {
			return err
		}
	}

	err = errors.New("invalid Terraform version, run `tfvm install --list` for a list of available versions")
	return err
}

// GetInstalledVersions returns a list of all installed Terraform versions.
func GetInstalledVersions(installPath string, extension string) ([]string, error) {
	var versions []string
	var err error = nil

	files, err := ioutil.ReadDir(installPath)
	if err != nil {
		return versions, err
	}
	for _, f := range files {
		v := strings.TrimPrefix(f.Name(), "terraform")
		v = strings.TrimSuffix(v, extension)
		versions = append(versions, v)
	}

	return versions, err
}

// IsInstalledVersions returns true if the specified Terraform version is installed.
func IsInstalledVersion(installPath string, extension string, version string) error {
	var err error = nil

	versions, err := GetInstalledVersions(installPath, extension)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if v == version {
			return err
		}
	}

	err = errors.New("invalid Terraform version, run `tfvm list` for a list of installed versions")
	return err
}
