package update

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/atrox/homedir"
	"github.com/hashicorp/go-version"
	"github.com/mattn/go-isatty"

	"github.com/confluentinc/cli/log"
)

type Client struct {
	repository    Repository
	lastCheckFile string
	logger        *log.Logger
}

// NewUpdateClient returns a client for updating CLI binaries
func NewUpdateClient(repo Repository, lastCheckFile string, logger *log.Logger) *Client {
	return &Client{
		repository:    repo,
		lastCheckFile: lastCheckFile,
		logger:        logger,
	}
}

// CheckForUpdates checks for new versions in the repo
func (c *Client) CheckForUpdates(name string, currentVersion string) (updateAvailable bool, latestVersion string, err error) {
	availableVersions, err := c.repository.GetAvailableVersions(name)
	if err != nil {
		return false, "", err
	}

	mostRecentVersion := availableVersions[len(availableVersions)-1]

	currVersion, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, "", fmt.Errorf("unable to parse %s version %s - %s", name, currentVersion, err)
	}

	if currVersion.LessThan(mostRecentVersion) {
		return true, mostRecentVersion.String(), nil
	}

	return false, currentVersion, nil
}

// PromptToDownload displays an interactive CLI prompt to download the latest version
func (c *Client) PromptToDownload(name, currVersion, latestVersion string, confirm bool) bool {
	if confirm && !isatty.IsTerminal(os.Stdout.Fd()) {
		c.logger.Warn("disable confirm as stdout is not a tty")
		confirm = false
	}

	c.logger.Printf("New version of %s is available", name)
	c.logger.Printf("Current Version: %s", currVersion)
	c.logger.Printf("Latest Version:  %s", latestVersion)

	if confirm {
		for {
			fmt.Print("Do you want to download and install this update? (y/n): ")

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')

			choice := string([]byte(input)[0])
			switch choice {
			case "y":
				return true
			case "n":
				return false
			default:
				fmt.Printf("%s is not a valid choice", choice)
				continue
			}
		}
	}

	return false
}

// UpdateBinary replaces the named binary at path with the desired version
func (c *Client) UpdateBinary(name, version, path string) error {
	downloadDir, err := ioutil.TempDir("", name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(downloadDir)

	c.logger.Printf("Downloading %s version %s...", name, version)
	startTime := time.Now()

	newBin, bytes, err := c.repository.DownloadVersion(name, version, downloadDir)
	if err != nil {
		return err
	}

	mb := float64(bytes) / 1024.0 / 1024.0
	timeSpent := time.Since(startTime).Seconds()
	c.logger.Printf("Done. Downloaded %.2f MB in %.0f seconds. (%.2f MB/s)", mb, timeSpent, mb/timeSpent)

	err = copyFile(newBin, path)
	if err != nil {
		return err
	}

	if err := os.Chmod(path, 0755); err != nil {
		return err
	}

	return nil
}


func (c *Client) TouchUpdateCheckFile() error {
	updateFile, err := homedir.Expand(LastCheckFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(updateFile); os.IsNotExist(err) {
		if f, err := os.Create(updateFile); err != nil {
			return err
		} else {
			f.Close()
		}
	} else if err := os.Chtimes(updateFile, time.Now(), time.Now()); err != nil {
		return err
	}
	return nil
}

// copyFile copies from src to dst until either EOF is reached
// on src or an error occurs. It verifies src exists and removes
// the dst if it exists.
func copyFile(src, dst string) error {
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	if cleanSrc == cleanDst {
		return nil
	}
	sf, err := os.Open(cleanSrc)
	if err != nil {
		return err
	}
	defer sf.Close()
	if err := os.Remove(cleanDst); err != nil && !os.IsNotExist(err) {
		return err
	}
	df, err := os.Create(cleanDst)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}
