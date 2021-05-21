package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getInstalledVersions() ([]string, error) {
	var versions []string
	var err error = nil

	files, err := ioutil.ReadDir(installPath)
	if err != nil {
		return versions, err
	}
	for _, f := range files {
		v := strings.TrimPrefix(f.Name(), "terraform")
		versions = append(versions, v)
	}

	return versions, err
}

func isInstalledVersion(version string) (bool, error) {
	var err error = nil

	versions, err := getInstalledVersions()
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if v == version {
			return true, err
		}
	}
	return false, err
}

func getAvailableVersions() ([]string, error) {
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
			versions = append(versions, version)
		}
	})
	return versions, err
}

func getLatestVersion() (string, error) {
	var err error = nil

	versions, err := getAvailableVersions()
	if err != nil {
		return "", err
	}
	return versions[0], err
}

func isAvailableVersion(version string) (bool, error) {
	var err error = nil

	versions, err := getAvailableVersions()
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if v == version {
			return true, err
		}
	}
	return false, err
}

func getArchitecture() (string, error) {
	var arch string = ""
	var err error = nil

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
		default:
			err = errors.New("architecture could not be verified for installation")
		}
	} else {
		err = errors.New("operating system is not supported for installation")
	}

	return arch, err
}

func downloadArchive(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("/tmp/tfvm.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

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
