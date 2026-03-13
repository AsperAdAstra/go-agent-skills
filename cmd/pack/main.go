package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	skillPath := flag.String("path", ".", "Path to skill directory")
	output := flag.String("output", "", "Output zip file (default: skill-name.zip)")
	flag.Parse()

	result, err := Pack(*skillPath, *output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Packed: %s (%d files, %d bytes)\n", result.Filename, result.FileCount, result.Size)
}

type PackResult struct {
	Filename string
	FileCount int
	Size     int64
}

func Pack(skillPath, outputPath string) (PackResult, error) {
	result := PackResult{}

	// Get skill name
	skillName := filepath.Base(skillPath)
	if skillName == "." {
		skillName = "skill"
	}

	// Default output
	if outputPath == "" {
		outputPath = skillName + ".zip"
	}
	result.Filename = outputPath

	// Create zip file
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// Walk directory
	err := filepath.Walk(skillPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files except .env.example
		relPath, _ := filepath.Rel(skillPath, path)
		if relPath == "." {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip .git, node_modules, etc.
		skipDirs := []string{".git", "node_modules", ".npm", "__pycache__"}
		for _, skip := range skipDirs {
			if strings.HasPrefix(relPath, skip) || relPath == skip {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Create zip entry
		header := &zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate,
		}
		if info.IsDir() {
			header.Name += "/"
		}
		
		f, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			f.Write(data)
			result.FileCount++
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	// Write zip file
	writer.Close()
	err = os.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		return result, err
	}

	result.Size = int64(buf.Len())
	return result, nil
}
