package gh

// Response struct for https://api.github.com/advisories
type Response struct {
	Severity        string `json:"severity"`
	CVE             string `json:"cve_id"`
	Summary         string `json:"summary"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Vulnerabilities []struct {
		Package struct {
			Name string `json:"name"`
		} `json:"package"`
		VulnerableVersionRange string   `json:"vulnerable_version_range"`
		FirstPatchedVersion    string   `json:"first_patched_version"`
		VulnerableFunctions    []string `json:"vulnerable_functions"`
	} `json:"vulnerabilities"`
	References []string `json:"references"`
}
