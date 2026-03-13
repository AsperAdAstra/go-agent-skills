package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	skillPath := flag.String("path", ".", "Path to skill directory")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	flag.Parse()

	result := Validate(*skillPath)

	if *jsonOutput {
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	} else {
		PrintResult(result)
	}

	if !result.Valid {
		os.Exit(1)
	}
}

type ValidationResult struct {
	Valid     bool              `json:"valid"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Errors    []string         `json:"errors,omitempty"`
	Warnings  []string         `json:"warnings,omitempty"`
	Files     []string         `json:"files"`
	Metadata  SkillMetadata    `json:"metadata,omitempty"`
}

type SkillMetadata struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Homepage   string            `json:"homepage"`
	Requires   MetadataRequires  `json:"requires"`
}

type MetadataRequires struct {
	Bins []string `json:"bins"`
	Env  []string `json:"env"`
}

func Validate(path string) ValidationResult {
	result := ValidationResult{
		Path: path,
	}

	// Get skill name from directory
	result.Name = filepath.Base(path)

	// Check required files
	requiredFiles := []string{"SKILL.md"}
	for _, f := range requiredFiles {
		fp := filepath.Join(path, f)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required file: %s", f))
		} else {
			result.Files = append(result.Files, f)
		}
	}

	// Check for scripts directory
	scriptsDir := filepath.Join(path, "scripts")
	if info, err := os.Stat(scriptsDir); err == nil && info.IsDir() {
		result.Files = append(result.Files, "scripts/")
		
		// Check for executable scripts
		entries, _ := os.ReadDir(scriptsDir)
		for _, e := range entries {
			if !e.IsDir() {
				result.Files = append(result.Files, "scripts/"+e.Name())
			}
		}
	} else {
		result.Warnings = append(result.Warnings, "No scripts/ directory found")
	}

	// Parse SKILL.md for metadata
	skillFile := filepath.Join(path, "SKILL.md")
	if data, err := os.ReadFile(skillFile); err == nil {
		metadata := parseMetadata(string(data))
		if metadata.Name != "" {
			result.Metadata = metadata
		}
		
		// Check for required metadata fields
		if metadata.Description == "" {
			result.Warnings = append(result.Warnings, "No description in SKILL.md")
		}
	}

	// Check for config or additional files
	additionalFiles := []string{"README.md", "config.json", ".env.example"}
	for _, f := range additionalFiles {
		fp := filepath.Join(path, f)
		if _, err := os.Stat(fp); err == nil {
			result.Files = append(result.Files, f)
		}
	}

	// Determine validity
	result.Valid = len(result.Errors) == 0

	return result
}

func parseMetadata(content string) SkillMetadata {
	metadata := SkillMetadata{}

	// Parse YAML frontmatter
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			yamlContent := parts[1]
			
			// Simple YAML parsing
			lines := strings.Split(yamlContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "name:") {
					metadata.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				}
				if strings.HasPrefix(line, "description:") {
					metadata.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
				}
				if strings.HasPrefix(line, "homepage:") {
					metadata.Homepage = strings.TrimSpace(strings.TrimPrefix(line, "homepage:"))
				}
			}
		}
	}

	return metadata
}

func PrintResult(r ValidationResult) {
	fmt.Printf("=== Skill Validation: %s ===\n", r.Name)
	fmt.Printf("Path: %s\n", r.Path)
	fmt.Printf("Valid: %v\n", r.Valid)
	
	if len(r.Errors) > 0 {
		fmt.Println("\n❌ Errors:")
		for _, e := range r.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}
	
	if len(r.Warnings) > 0 {
		fmt.Println("\n⚠️  Warnings:")
		for _, w := range r.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}
	
	if len(r.Files) > 0 {
		fmt.Println("\n📁 Files:")
		for _, f := range r.Files {
			fmt.Printf("  - %s\n", f)
		}
	}
	
	if r.Metadata.Name != "" {
		fmt.Println("\n📝 Metadata:")
		fmt.Printf("  Name: %s\n", r.Metadata.Name)
		fmt.Printf("  Description: %s\n", r.Metadata.Description)
	}
}
