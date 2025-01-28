package main

import (
	"flag"
	"fmt"
	"github.com/mlw157/scout/internal/detectors/filesystem"
	"github.com/mlw157/scout/internal/engine"
	"github.com/mlw157/scout/internal/exporters/jsonexporter"
	"log"
	"strings"
)

// just to test
func main() {

	ecosystemsFlag := flag.String("ecosystems", "", "Comma-separated list of ecosystems to scan (e.g., go,pip,maven)")
	excludeDirsFlag := flag.String("exclude", "", "Comma-separated list of directory and file names to exclude (e.g., node_modules,.git,requirements-dev.txt)")
	exportFlag := flag.Bool("export", false, "Export results to a file (default is no export)")
	tokenFlag := flag.String("token", "", "GitHub token for authenticated API requests (optional)")
	sequentialFlag := flag.Bool("sequential", false, "Processes each file individually without concurrent execution (not recommended)")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Please provide a root directory to scan")
	}

	rootDir := args[0]

	// ecosystems flag
	var ecosystems []string

	if *ecosystemsFlag != "" {
		ecosystems = strings.Split(*ecosystemsFlag, ",")
	} else {
		// default ecosystems
		ecosystems = []string{"go", "maven", "pip", "npm", "composer"}
	}

	// exclude directories flag
	var excludeDirs []string

	if *excludeDirsFlag != "" {
		excludeDirs = strings.Split(*excludeDirsFlag, ",")
	} else {
		excludeDirs = []string{}
	}

	log.Printf("Path to scan: %s\n", rootDir)
	log.Printf("Ecosystems to scan: %v\n", ecosystems)
	log.Printf("Excluded directories: %v\n", excludeDirs)

	detector := filesystem.NewFSDetector()

	config := engine.Config{
		Ecosystems:     ecosystems,
		ExcludeFiles:   excludeDirs,
		Token:          *tokenFlag,
		SequentialMode: *sequentialFlag,
	}

	// if export flag is set, create a exporter
	// todo make multiple export types, other than json
	if *exportFlag {
		config.Exporter = jsonexporter.NewJSONExporter("scout_report.json")
	}

	scanEngine := engine.NewEngine(detector, config)

	scanResults, err := scanEngine.Scan(rootDir)
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	log.Printf("Scan results for directory: %s\n\n", rootDir)

	for _, result := range scanResults {
		log.Println("File: " + result.SourceFile)
		log.Printf("Found %d vulnerabilities in %d packages\n\n", len(result.Vulnerabilities), len(result.Dependencies))

		if len(result.Vulnerabilities) > 0 {
			log.Println("Vulnerabilities found:")
			fmt.Println()
			for _, vulnerability := range result.Vulnerabilities {
				log.Printf("Package: %s@%s\n", vulnerability.Dependency.Name, vulnerability.Dependency.Version)
				log.Printf("CVE: %s\n", vulnerability.CVE)
				log.Printf("Severity: %s\n", vulnerability.Severity)
				log.Printf("Summary: %s\n", vulnerability.Summary)
				log.Printf("Upgrade to version %s in order to fix\n", vulnerability.FirstPatchedVersion)
				log.Println()

			}
		}
	}

}
