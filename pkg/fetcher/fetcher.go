package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zinrai/freebsd-archive-combiner/pkg/config"
)

func FetchSplitFiles(cfg *config.Config, component *config.Component) error {
	combinedFilePath := config.GetCombinedFilePath(cfg, component)

	if fileExists(combinedFilePath) {
		fmt.Printf("Combined file already exists: %s\n", combinedFilePath)
		return nil
	}

	fetchDir := config.GetFetchDir(cfg, component)

	baseURL := fmt.Sprintf("%s/%s/%s/%s",
		cfg.ArchiveURL,
		cfg.Architecture,
		cfg.Version,
		component.Directory)

	fmt.Printf("Fetching files for component: %s\n", component.Directory)
	fmt.Printf("URL: %s\n", baseURL)
	fmt.Printf("Destination: %s\n", fetchDir)

	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < len(alphabet); i++ {
		for j := 0; j < len(alphabet); j++ {
			suffix := string(alphabet[i]) + string(alphabet[j])
			filename := component.FilePrefix + "." + suffix
			fileURL := fmt.Sprintf("%s/%s", baseURL, filename)
			destPath := filepath.Join(fetchDir, filename)

			if fileExists(destPath) {
				fmt.Printf("  Using existing file: %s\n", filename)
				continue
			}

			exists, err := checkFileExists(fileURL)
			if err != nil {
				return fmt.Errorf("error checking file existence: %w", err)
			}

			if !exists {
				if suffix == "aa" {
					return fmt.Errorf("first file %s does not exist at %s", filename, fileURL)
				}
				fmt.Printf("  No more files to download for %s\n", component.Directory)
				return nil
			}

			fmt.Printf("  Downloading: %s\n", filename)
			if err := downloadFile(destPath, fileURL); err != nil {
				return fmt.Errorf("failed to download %s: %w", filename, err)
			}
		}
	}

	return nil
}

func checkFileExists(url string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func downloadFile(destPath, url string) error {
	tmpFile := destPath + ".download"

	os.Remove(tmpFile)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, destPath)
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
