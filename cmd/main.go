package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zinrai/freebsd-archive-combiner/pkg/combiner"
	"github.com/zinrai/freebsd-archive-combiner/pkg/config"
	"github.com/zinrai/freebsd-archive-combiner/pkg/fetcher"
)

func main() {
	configFile := flag.String("c", "", "Path to YAML configuration file")
	showHelp := flag.Bool("h", false, "Show help information")

	flag.Parse()

	if *showHelp || *configFile == "" {
		printUsage()
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if err := config.EnsureOutputDirs(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing FreeBSD %s for %s architecture\n", cfg.Version, cfg.Architecture)

	for i := range cfg.Components {
		component := &cfg.Components[i]

		fmt.Printf("\n=== Processing component: %s ===\n", component.Directory)

		err := fetcher.FetchSplitFiles(cfg, component)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching files for %s: %v\n", component.Directory, err)
			os.Exit(1)
		}

		err = combiner.CombineSplitFiles(cfg, component)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error combining files for %s: %v\n", component.Directory, err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nAll components processed successfully.\n")
	fmt.Printf("Combined files are available in: output/%s/%s/combine/\n", cfg.Version, cfg.Architecture)
}

func printUsage() {
	fmt.Println("freebsd-archive-combiner - Downloads and combines FreeBSD split files from archive")
	fmt.Println("\nUsage:")
	fmt.Println("  freebsd-archive-combiner -c CONFIG_FILE.yaml")
	fmt.Println("\nOptions:")
	fmt.Println("  -c CONFIG_FILE.yaml   Path to YAML configuration file")
	fmt.Println("  -h                    Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  freebsd-archive-combiner -c 8.4-RELEASE.yaml")
}
