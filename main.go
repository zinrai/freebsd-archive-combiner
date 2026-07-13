package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	configFile := flag.String("config", "", "Path to YAML configuration file")
	outputDir := flag.String("output", "output", "Root directory for downloaded and combined files")
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("h", false, "Show help information")

	flag.Parse()

	if *showVersion {
		fmt.Printf("freebsd-archive-combiner %s (commit %s, built %s)\n", version, commit, date)
		os.Exit(0)
	}

	if *showHelp || *configFile == "" {
		printUsage()
		os.Exit(0)
	}

	cfg, err := LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if err := EnsureOutputDirs(*outputDir, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing FreeBSD %s for %s architecture\n", cfg.Version, cfg.Architecture)

	for i := range cfg.Components {
		component := &cfg.Components[i]

		fmt.Printf("\n=== Processing component: %s ===\n", component.Directory)

		err := FetchSplitFiles(*outputDir, cfg, component)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching files for %s: %v\n", component.Directory, err)
			os.Exit(1)
		}

		err = CombineSplitFiles(*outputDir, cfg, component)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error combining files for %s: %v\n", component.Directory, err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nAll components processed successfully.\n")
	fmt.Printf("Combined files are available in: %s\n", GetCombineDir(*outputDir, cfg))
}

func printUsage() {
	fmt.Println("freebsd-archive-combiner - Downloads and combines FreeBSD split files from archive")
	fmt.Println("\nUsage:")
	fmt.Println("  freebsd-archive-combiner -config CONFIG_FILE.yaml")
	fmt.Println("\nOptions:")
	fmt.Println("  -config CONFIG_FILE.yaml   Path to YAML configuration file")
	fmt.Println("  -output DIR                Root directory for output (default: output)")
	fmt.Println("  -version                   Show version information")
	fmt.Println("  -h                         Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  freebsd-archive-combiner -config 8.4-RELEASE.yaml")
}
