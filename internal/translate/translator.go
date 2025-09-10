package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Provider represents a translation service provider
type Provider interface {
	Translate(text, sourceLang, targetLang string) (string, error)
}

// OpenAIProvider implements translation using OpenAI API
type OpenAIProvider struct {
	APIKey string
	Model  string
}

// AnthropicProvider implements translation using Anthropic API
type AnthropicProvider struct {
	APIKey string
	Model  string
}

// OllamaProvider implements translation using local Ollama
type OllamaProvider struct {
	BaseURL string
	Model   string
}

// NewProvider creates a new translation provider
func NewProvider(providerType, apiKey, model string) (Provider, error) {
	switch strings.ToLower(providerType) {
	case "openai":
		if model == "" {
			model = "gpt-4"
		}
		return &OpenAIProvider{
			APIKey: apiKey,
			Model:  model,
		}, nil
	case "anthropic":
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
		return &AnthropicProvider{
			APIKey: apiKey,
			Model:  model,
		}, nil
	case "ollama":
		if model == "" {
			model = getEnvDefault("OLLAMA_MODEL", "llama3:8b")
		}
		return &OllamaProvider{
			BaseURL: getEnvDefault("OLLAMA_URL", "http://localhost:11434"),
			Model:   model,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}

// TranslateTemplate translates a template file while preserving all structure
func TranslateTemplate(content, sourceLang, targetLang, providerType, model string) (string, error) {
	provider, err := NewProvider(providerType, getAPIKey(providerType), model)
	if err != nil {
		return "", err
	}

	// Use direct translation approach with structured prompt
	return TranslateTemplateDirect(content, sourceLang, targetLang, provider)
}

// TranslateTemplateDirect uses a single LLM call with a structured prompt to translate the entire document
func TranslateTemplateDirect(content, sourceLang, targetLang string, provider Provider) (string, error) {
	// Create structured prompt based on translate-prompt.md
	prompt := fmt.Sprintf(`You are a professional translator. Your task is to translate a markdown template document while preserving its exact structure and format.

**Instructions:**
1. Translate the YAML metadata fields: `+"`name`"+` and `+"`comment`"+` values to %s
2. Translate all content after the `+"`---`"+` separator to %s
3. Keep all other YAML fields unchanged (acronym, satisfies, dates, etc.)
4. Preserve all markdown formatting, template variables (like {{.Name}}), and document structure
5. Return only the translated content with no additional comments or explanations
6. Maintain the exact same line breaks and spacing as the original

**Input document:**
%s

**Target language:** %s`, targetLang, targetLang, content, targetLang)

	// Call the provider with the structured prompt
	result, err := provider.Translate(prompt, sourceLang, targetLang)
	if err != nil {
		return "", errors.Wrapf(err, "failed to translate template using direct approach")
	}

	// Light cleanup - only remove obvious artifacts while preserving structure
	cleaned := cleanDirectTranslationResponse(result, content)

	return cleaned, nil
}

// cleanDirectTranslationResponse performs minimal cleaning while preserving document structure
func cleanDirectTranslationResponse(response, original string) string {
	// Remove common prefixes that might appear before the document
	prefixPatterns := []string{
		`(?i)^\s*here\s+is\s+the\s+translated\s+document:?\s*\n?`,
		`(?i)^\s*translated\s+document:?\s*\n?`,
		`(?i)^\s*result:?\s*\n?`,
		`(?i)^\s*translation:?\s*\n?`,
	}

	result := response
	for _, pattern := range prefixPatterns {
		result = regexp.MustCompile(pattern).ReplaceAllString(result, "")
	}

	// If result doesn't start with expected YAML frontmatter or content, try to extract it
	lines := strings.Split(result, "\n")
	startIndex := -1

	// Look for the start of meaningful content
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if line looks like YAML field or --- delimiter or markdown header
		if trimmed == "---" ||
			regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*:\s*`).MatchString(trimmed) ||
			strings.HasPrefix(trimmed, "#") {
			startIndex = i
			break
		}
	}

	if startIndex > 0 {
		result = strings.Join(lines[startIndex:], "\n")
	}

	return strings.TrimSpace(result)
}

// TranslateDocument translates a compliance document while preserving structure
func TranslateDocument(content, sourceLang, targetLang, providerType, model string) (string, error) {
	provider, err := NewProvider(providerType, getAPIKey(providerType), model)
	if err != nil {
		return "", err
	}

	// Extract metadata and content sections
	sections := extractSections(content)

	var translatedSections []string

	for _, section := range sections {
		if shouldTranslate(section) {
			translated, err := provider.Translate(section, sourceLang, targetLang)
			if err != nil {
				return "", errors.Wrapf(err, "failed to translate section")
			}
			translatedSections = append(translatedSections, translated)
		} else {
			// Keep metadata, tables, and code blocks as-is
			translatedSections = append(translatedSections, section)
		}
	}

	return strings.Join(translatedSections, "\n"), nil
}

// OpenAI API implementation
func (p *OpenAIProvider) Translate(text, sourceLang, targetLang string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional compliance document translator. Translate the following text from %s to %s exactly as written, preserving all formatting, markdown syntax, YAML frontmatter, and technical terms. Do not add any comments, explanations, or additional content. Return only the translated text.

%s`, sourceLang, targetLang, text)

	requestBody := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  4000,
		"temperature": 0.1, // Low temperature for consistency
	}

	return p.makeRequest("https://api.openai.com/v1/chat/completions", requestBody)
}

// Anthropic API implementation
func (p *AnthropicProvider) Translate(text, sourceLang, targetLang string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional compliance document translator. Translate the following text from %s to %s exactly as written, preserving all formatting, markdown syntax, YAML frontmatter, and technical terms. Do not add any comments, explanations, or additional content. Return only the translated text.

%s`, sourceLang, targetLang, text)

	requestBody := map[string]interface{}{
		"model":      p.Model,
		"max_tokens": 4000,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	return p.makeRequest("https://api.anthropic.com/v1/messages", requestBody)
}

// Ollama API implementation
func (p *OllamaProvider) Translate(text, sourceLang, targetLang string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional compliance document translator. Translate the following text from %s to %s exactly as written, preserving all formatting, markdown syntax, YAML frontmatter, and technical terms. Do not add any comments, explanations, or additional content. Return only the translated text.

%s`, sourceLang, targetLang, text)

	requestBody := map[string]interface{}{
		"model":  p.Model,
		"prompt": prompt,
		"stream": false,
	}

	url := fmt.Sprintf("%s/api/generate", p.BaseURL)
	return p.makeRequest(url, requestBody)
}

func (p *OpenAIProvider) makeRequest(url string, requestBody map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed: %s", string(body))
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Extract content from OpenAI response
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return strings.TrimSpace(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

func (p *AnthropicProvider) makeRequest(url string, requestBody map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed: %s", string(body))
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Extract content from Anthropic response
	if content, ok := response["content"].([]interface{}); ok && len(content) > 0 {
		if item, ok := content[0].(map[string]interface{}); ok {
			if text, ok := item["text"].(string); ok {
				return strings.TrimSpace(text), nil
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

func (p *OllamaProvider) makeRequest(url string, requestBody map[string]interface{}) (string, error) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second} // Longer timeout for local models
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed: %s", string(body))
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Extract content from Ollama response
	if response["response"] != nil {
		return strings.TrimSpace(response["response"].(string)), nil
	}

	return "", fmt.Errorf("unexpected response format")
}

// extractSections splits document into translatable sections
func extractSections(content string) []string {
	// Split by paragraphs but keep structure
	lines := strings.Split(content, "\n")
	var sections []string
	var currentSection []string

	for _, line := range lines {
		if line == "" {
			if len(currentSection) > 0 {
				sections = append(sections, strings.Join(currentSection, "\n"))
				currentSection = []string{}
			}
			sections = append(sections, line)
		} else {
			currentSection = append(currentSection, line)
		}
	}

	if len(currentSection) > 0 {
		sections = append(sections, strings.Join(currentSection, "\n"))
	}

	return sections
}

// shouldTranslate determines if a section should be translated
func shouldTranslate(section string) bool {
	section = strings.TrimSpace(section)

	// Skip empty sections
	if section == "" {
		return false
	}

	// Skip YAML frontmatter
	if strings.HasPrefix(section, "---") || strings.HasSuffix(section, "---") {
		return false
	}

	// Skip LaTeX commands
	if strings.HasPrefix(section, "%") {
		return false
	}

	// Skip tables (markdown tables)
	if strings.Contains(section, "|") && (strings.Contains(section, "---") || strings.Contains(section, "Table:")) {
		return false
	}

	// Skip code blocks
	if strings.HasPrefix(section, "```") || strings.HasSuffix(section, "```") {
		return false
	}

	// Skip header-includes and other metadata
	if strings.Contains(section, "header-includes:") || strings.Contains(section, "\\usepackage") {
		return false
	}

	return true
}

// getAPIKey retrieves API key from environment
func getAPIKey(provider string) string {
	switch strings.ToLower(provider) {
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	default:
		return ""
	}
}

// getEnvDefault returns environment variable or default value
func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
