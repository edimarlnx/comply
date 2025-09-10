package path

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
)

// File wraps an os.FileInfo as well as the absolute path to the underlying file.
type File struct {
	FullPath string
	Info     os.FileInfo
}

// Standards lists all standard files.
var Standards = func() ([]File, error) {
	return filesFor("standards", "yml")
}

// Narratives lists all narrative files.
var Narratives = func() ([]File, error) {
	return filesFor("narratives", "md")
}

// Policies lists all policy files.
var Policies = func() ([]File, error) {
	return filesFor("policies", "md")
}

// Procedures lists all procedure files.
var Procedures = func() ([]File, error) {
	return filesFor("procedures", "md")
}

func filesFor(name, extension string) ([]File, error) {
	var filtered []File
	files, err := ioutil.ReadDir(filepath.Join(".", name))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load files for: "+name)
	}

	// Track original files and their translations
	originalFiles := make(map[string]File)
	translatedFiles := make(map[string][]File)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "."+extension) || strings.HasPrefix(strings.ToUpper(f.Name()), "README") {
			continue
		}

		abs, err := filepath.Abs(filepath.Join(".", name, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "unable to load file: "+f.Name())
		}

		file := File{abs, f}

		// Check if this is a translated file
		if lang := extractLanguageFromFilename(f.Name(), extension); lang != "" {
			// This is a translated file
			baseName := getBaseFilename(f.Name(), lang, extension)
			translatedFiles[baseName] = append(translatedFiles[baseName], file)
		} else {
			// This is an original file
			originalFiles[f.Name()] = file
		}
	}

	// First add all original files
	for _, file := range originalFiles {
		filtered = append(filtered, file)
	}

	// Then add translated files if translation is enabled
	cfg := config.Config()
	if cfg.Translation != nil && cfg.Translation.Enabled {
		for _, translations := range translatedFiles {
			// Only add translations for configured languages
			for _, translatedFile := range translations {
				lang := extractLanguageFromFilename(translatedFile.Info.Name(), extension)
				if contains(cfg.Translation.Languages, lang) {
					filtered = append(filtered, translatedFile)
				}
			}
		}
	}

	return filtered, nil
}

// extractLanguageFromFilename extracts language code from translated filename
// e.g., "policy.pt-BR.md" -> "pt-BR"
func extractLanguageFromFilename(filename, extension string) string {
	parts := strings.Split(filename, ".")
	if len(parts) >= 3 {
		// Check if second-to-last part looks like a language code
		langPart := parts[len(parts)-2]
		if len(langPart) >= 2 && (strings.Contains(langPart, "-") || len(langPart) == 2) {
			return langPart
		}
	}
	return ""
}

// getBaseFilename gets original filename from translated filename
// e.g., "policy.pt-BR.md" -> "policy.md"
func getBaseFilename(filename, lang, extension string) string {
	return strings.Replace(filename, "."+lang+"."+extension, "."+extension, 1)
}

// contains checks if a string slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
