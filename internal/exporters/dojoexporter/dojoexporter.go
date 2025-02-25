package dojoexporter

import (
	"encoding/json"
	"fmt"
	"github.com/mlw157/scout/internal/models"
	"log"
	"os"
)

type DojoExporter struct {
	OutputFile string
}

// NewDojoExporter creates a new DojoExporter
func NewDojoExporter(outputFile string) *DojoExporter {
	return &DojoExporter{OutputFile: outputFile}
}

// DojoFinding represents the minimal required fields for DefectDojo
type DojoFinding struct {
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	// Optional field
	CVE string `json:"cve,omitempty"`
}

// Export converts scan results to minimal DefectDojo format and saves to file
func (d *DojoExporter) Export(results []*models.ScanResult) error {
	file, err := os.Create(d.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", d.OutputFile, err)
	}
	defer file.Close()

	var dojoFindings []DojoFinding

	for _, result := range results {
		for _, vulnerability := range result.Vulnerabilities {
			// Map Scout severity to DefectDojo severity
			severity := mapSeverity(vulnerability.Severity)

			// Create a more comprehensive description that includes file path and affected component
			enhancedDescription := fmt.Sprintf("%s\n\nAffected File: %s\nPackage: %s@%s\n",
				vulnerability.Description,
				result.SourceFile,
				vulnerability.Dependency.Name,
				vulnerability.Dependency.Version)

			// Add remediation info if available
			if vulnerability.FirstPatchedVersion != "" {
				enhancedDescription += fmt.Sprintf("Remediation: Update to version %s or later\n",
					vulnerability.FirstPatchedVersion)
			}

			dojoFinding := DojoFinding{
				Title:       vulnerability.Summary,
				Severity:    severity,
				Description: enhancedDescription,
				CVE:         vulnerability.CVE,
			}

			dojoFindings = append(dojoFindings, dojoFinding)
		}
	}

	// Create the final Dojo report structure with just findings
	dojoReport := map[string]interface{}{
		"findings": dojoFindings,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // pretty-print JSON
	if err := encoder.Encode(dojoReport); err != nil {
		return fmt.Errorf("failed to encode results to DefectDojo format: %v", err)
	}

	log.Printf("Vulnerabilities exported to %s in DefectDojo format\n", d.OutputFile)
	return nil
}

// mapSeverity maps Scout severity levels to DefectDojo severity levels
func mapSeverity(scoutSeverity string) string {
	switch scoutSeverity {
	case "critical":
		return "Critical"
	case "high":
		return "High"
	case "medium":
		return "Medium"
	case "low":
		return "Low"
	default:
		return "Info" // Default to Info if unknown
	}
}
