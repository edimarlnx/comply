package render

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/translate"
)

// TranslateTemplates translates template files based on path or config
func TranslateTemplates(path, provider string) error {
	cfg := config.Config()

	// Get languages from config or default
	var languages []string
	var providerToUse string
	var modelToUse string

	if cfg.Translation != nil && cfg.Translation.Enabled {
		languages = cfg.Translation.Languages
		if provider == "" {
			providerToUse = cfg.Translation.Provider
		} else {
			providerToUse = provider
		}
		modelToUse = cfg.Translation.Model
	} else {
		return fmt.Errorf("translation not enabled in comply.yml")
	}

	if len(languages) == 0 {
		return fmt.Errorf("no languages configured for translation")
	}

	if providerToUse == "" {
		return fmt.Errorf("no translation provider specified")
	}

	// Determine files to translate
	var filesToTranslate []string
	var err error

	if path != "" {
		// Translate specific path
		filesToTranslate, err = getFilesToTranslate(path)
		if err != nil {
			return errors.Wrapf(err, "unable to get files from path: %s", path)
		}
	} else {
		// Translate all templates
		filesToTranslate, err = getAllTemplateFiles()
		if err != nil {
			return errors.Wrap(err, "unable to get all template files")
		}
	}

	if len(filesToTranslate) == 0 {
		fmt.Println("No files to translate")
		return nil
	}

	fmt.Printf("Translating %d files to languages: %s\n", len(filesToTranslate), strings.Join(languages, ", "))

	// Translate each file to each language
	for _, file := range filesToTranslate {
		for _, lang := range languages {
			err = translateSingleTemplate(file, lang, providerToUse, modelToUse)
			if err != nil {
				return errors.Wrapf(err, "failed to translate %s to %s", file, lang)
			}
		}
	}

	fmt.Println("Template translation completed successfully")
	return nil
}

// getFilesToTranslate returns list of files to translate from given path
func getFilesToTranslate(path string) ([]string, error) {
	var files []string

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "path does not exist: %s", path)
	}

	if info.IsDir() {
		// Directory - find all .md files
		err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(p, ".md") && isTemplateFile(p) {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "error walking directory: %s", path)
		}
	} else {
		// Single file
		if strings.HasSuffix(path, ".md") && isTemplateFile(path) {
			files = append(files, path)
		} else {
			return nil, fmt.Errorf("file is not a translatable template: %s", path)
		}
	}

	return files, nil
}

// getAllTemplateFiles returns all template files in the project
func getAllTemplateFiles() ([]string, error) {
	var files []string

	// Standard directories containing templates
	dirs := []string{"policies", "procedures", "narratives"}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue // Skip if directory doesn't exist
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".md") && isTemplateFile(path) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "error walking directory: %s", dir)
		}
	}

	return files, nil
}

// isTemplateFile checks if file is a template (not already translated)
func isTemplateFile(path string) bool {
	// Skip files that are already translations (contain language code)
	filename := filepath.Base(path)

	// Skip README files and other non-template files
	if strings.HasPrefix(filename, "README") || strings.HasPrefix(filename, "TODO") {
		return false
	}

	// Check if filename contains language code pattern (.lang.)
	parts := strings.Split(filename, ".")
	if len(parts) >= 3 {
		// Check if second-to-last part looks like a language code
		langPart := parts[len(parts)-2]
		if len(langPart) >= 2 && (strings.Contains(langPart, "-") || len(langPart) == 2) {
			return false // This is already a translated file
		}
	}

	return true
}

// translateSingleTemplate translates a single template file
func translateSingleTemplate(filePath, targetLang, provider, model string) error {
	// Read original file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to read file: %s", filePath)
	}

	// Generate translated filename
	translatedPath := generateTranslatedPath(filePath, targetLang)

	// Check if translation already exists and is newer
	if shouldSkipTranslation(filePath, translatedPath) {
		fmt.Printf("Skipping %s (translation is up to date)\n", translatedPath)
		return nil
	}

	fmt.Printf("Translating %s to %s...\n", filePath, targetLang)

	// Translate content
	translatedContent, err := translate.TranslateTemplate(string(content), "en", targetLang, provider, model)
	if err != nil {
		return errors.Wrapf(err, "translation failed for file: %s", filePath)
	}

	// Write translated file
	err = ioutil.WriteFile(translatedPath, []byte(translatedContent), 0644)
	if err != nil {
		return errors.Wrapf(err, "unable to write translated file: %s", translatedPath)
	}

	fmt.Printf("Generated: %s\n", translatedPath)
	return nil
}

// generateTranslatedPath creates the translated filename
func generateTranslatedPath(originalPath, lang string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	translatedFilename := fmt.Sprintf("%s.%s%s", nameWithoutExt, lang, ext)
	return filepath.Join(dir, translatedFilename)
}

// shouldSkipTranslation checks if translation should be skipped
func shouldSkipTranslation(originalPath, translatedPath string) bool {
	// Check if translated file exists
	translatedInfo, err := os.Stat(translatedPath)
	if err != nil {
		return false // Translation doesn't exist, should translate
	}

	// Check if original file is newer than translation
	originalInfo, err := os.Stat(originalPath)
	if err != nil {
		return false // Can't stat original, should translate
	}

	// Skip if translation is newer than or equal to original
	return translatedInfo.ModTime().After(originalInfo.ModTime()) ||
		translatedInfo.ModTime().Equal(originalInfo.ModTime())
}
