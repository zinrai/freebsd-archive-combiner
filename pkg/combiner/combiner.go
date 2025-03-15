package combiner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zinrai/freebsd-archive-combiner/pkg/config"
)

func CombineSplitFiles(cfg *config.Config, component *config.Component) error {
	combinedFilePath := config.GetCombinedFilePath(cfg, component)

	if fileExists(combinedFilePath) {
		fmt.Printf("Combined file already exists: %s\n", combinedFilePath)
		return nil
	}

	fetchDir := config.GetFetchDir(cfg, component)

	fmt.Printf("Combining files for component: %s\n", component.Directory)
	fmt.Printf("Source directory: %s\n", fetchDir)
	fmt.Printf("Output file: %s\n", combinedFilePath)

	files, err := os.ReadDir(fetchDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", fetchDir, err)
	}

	prefix := component.FilePrefix + "."
	var splitFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			splitFiles = append(splitFiles, filepath.Join(fetchDir, file.Name()))
		}
	}

	if len(splitFiles) == 0 {
		return fmt.Errorf("no split files found for component %s in %s", component.Directory, fetchDir)
	}

	sort.Strings(splitFiles)

	fmt.Printf("Found %d split files to combine\n", len(splitFiles))

	err = combineFiles(splitFiles, combinedFilePath)
	if err != nil {
		return fmt.Errorf("failed to combine files: %w", err)
	}

	fmt.Printf("Successfully combined files to: %s\n", combinedFilePath)
	return nil
}

func combineFiles(files []string, outputPath string) error {
	tempFile := outputPath + ".tmp"

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	for i, file := range files {
		fmt.Printf("  Combining file %d/%d: %s\n", i+1, len(files), filepath.Base(file))

		in, err := os.Open(file)
		if err != nil {
			os.Remove(tempFile)
			return fmt.Errorf("failed to open input file %s: %w", file, err)
		}

		_, err = io.Copy(out, in)
		in.Close()

		if err != nil {
			os.Remove(tempFile)
			return fmt.Errorf("failed to copy from %s: %w", file, err)
		}
	}

	out.Close()

	if err := os.Rename(tempFile, outputPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
