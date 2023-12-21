package terraform

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	defaultVersion = "1.6.6"
	envTempDir     = "TF_TEMP_DIR"
	envVersion     = "TF_VERSION"
	envArch        = "TF_ARCH"
	envOS          = "TF_OS"
)

var (
	tempDir  string
	execPath string
)

func init() {
	var err error

	if os.Getenv(envTempDir) != "" {
		tempDir = os.Getenv(envTempDir)
	} else {
		tempDir, err = os.MkdirTemp(os.TempDir(), "tf")
		if err != nil {
			logrus.Fatal(err)
		}
	}
	logrus.Debugf("Using temp dir %s", tempDir)
	zipPath := filepath.Join(tempDir, "terraform.zip")
	if err := download(getDownloadURL(), zipPath); err != nil {
		logrus.Fatal(err)
	}
	if err := unzip(zipPath, tempDir); err != nil {
		logrus.Fatal(err)
	}
	execPath = filepath.Join(tempDir, "terraform")
	if err := os.Chmod(execPath, 0755); err != nil {
		logrus.Fatal(err)
	}
	output, err := exec.Command(execPath, "--version").CombinedOutput()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debugf("Terraform version: %s", output)
}

func getDownloadURL() string {
	version := defaultVersion
	if os.Getenv(envVersion) != "" {
		version = os.Getenv(envVersion)
	}
	arch := runtime.GOARCH
	if os.Getenv(envArch) != "" {
		arch = os.Getenv(envArch)
	}
	_os := runtime.GOOS
	if os.Getenv(envOS) != "" {
		_os = os.Getenv(envOS)
	}
	return fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
		version,
		version,
		_os,
		arch,
	)
}

func download(url, local string) error {
	logrus.Debugf("Downloading from %s to %s...", url, local)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(local)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func unzip(src, dest string) error {
	logrus.Debugf("Unzipping from %s to %s...", src, dest)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
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

func runTFWithEnv(args []string) (string, error) {
	cmd := exec.Command(execPath, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TF_IN_AUTOMATION=true")

	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	logrus.Debugf("Calling %s %v", execPath, args)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Error running terraform: %v", err)
		switch err.(type) {
		case *exec.ExitError:
			logrus.Debug(cmd.ProcessState.String())
			logrus.Debug(errBuffer.String())
			logrus.Debug(outBuffer.String())
		}
		return "", err
	}
	return outBuffer.String(), nil
}
