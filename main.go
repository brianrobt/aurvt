package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version = "1.0.0"
	Commit  = "development" // This can be set via build flags
	Date    = "unknown"     // This can be set via build flags
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type GitHubTag struct {
	Name string `json:"name"`
}

type PKGBUILDInfo struct {
	Name    string
	Version string
	URL     string
	Source  []string
}

// Variable substitution map
type Variables map[string]string

func main() {
	var rootCmd = &cobra.Command{
		Use:     "aurvt [package-directory]",
		Short:   "Check for newer versions of AUR packages on GitHub",
		Long:    `A CLI tool to check if newer versions are available for AUR packages hosted on GitHub.`,
		Version: Version,
		Args:    cobra.ExactArgs(1),
		Run:     runVersionCheck,
	}

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("aurvt version %s\n", Version)
			fmt.Printf("Commit: %s\n", Commit)
			fmt.Printf("Build Date: %s\n", Date)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runVersionCheck(cmd *cobra.Command, args []string) {
	packageDir := args[0]
	pkgbuildPath := fmt.Sprintf("%s/PKGBUILD", packageDir)

	// Check if package is in .aurvtignore file
	packageName := filepath.Base(packageDir)
	if isPackageIgnored(packageName) {
		fmt.Printf("📋 Package '%s' is ignored (found in .aurvtignore)\n", packageName)
		return
	}

	// Parse PKGBUILD
	pkgInfo, err := parsePKGBUILD(pkgbuildPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PKGBUILD: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Package: %s\n", pkgInfo.Name)
	fmt.Printf("Current version: %s\n", pkgInfo.Version)
	fmt.Printf("Repository URL: %s\n", pkgInfo.URL)

	if len(pkgInfo.Source) > 0 {
		fmt.Printf("Source URLs:\n")
		for i, source := range pkgInfo.Source {
			fmt.Printf("  [%d] %s\n", i+1, source)
		}
	}

	// Check if it's a GitHub repository
	if !strings.Contains(pkgInfo.URL, "github.com") {
		fmt.Println("❌ Not a GitHub repository, version checking not supported")
		return
	}

	// Get latest version from GitHub
	latestVersion, err := getLatestGitHubVersion(pkgInfo.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching latest version: %v\n", err)
		os.Exit(1)
	}

	// Check if source URLs use archive format and suggest tags endpoint
	checkAndSuggestURLFormat(pkgInfo.Source, pkgInfo.URL)

	// https://github.com/pyhunspell/pyhunspell/tags
	// https://github.com/pyhunspell/pyhunspell/archive/refs/tags/0.5.5.tar.gz
	if latestVersion == "" {
		fmt.Println("❌ Could not fetch latest version")
		return
	}

	fmt.Printf("Latest version: %s\n", latestVersion)

	// Compare versions
	if latestVersion == pkgInfo.Version {
		fmt.Println("✅ Package is up to date")
	} else {
		fmt.Printf("🔄 New version available: %s → %s\n", pkgInfo.Version, latestVersion)
	}
}

func parsePKGBUILD(filepath string) (*PKGBUILDInfo, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PKGBUILD: %w", err)
	}

	text := string(content)

	// Parse all variables first
	variables := parseVariables(text)

	// Extract pkgname with variable substitution
	pkgname := extractValue(text, `pkgname\s*=\s*(.+)`)
	if pkgname == "" {
		return nil, fmt.Errorf("could not find pkgname in PKGBUILD")
	}
	pkgname = substituteVariables(cleanValue(pkgname), variables)

	// Extract pkgver
	pkgver := extractValue(text, `pkgver\s*=\s*(.+)`)
	if pkgver == "" {
		return nil, fmt.Errorf("could not find pkgver in PKGBUILD")
	}
	pkgver = cleanValue(pkgver)

	// Extract url with variable substitution
	url := extractValue(text, `url\s*=\s*(.+)`)
	if url == "" {
		return nil, fmt.Errorf("could not find url in PKGBUILD")
	}
	url = substituteVariables(cleanValue(url), variables)

	// Extract source array with variable substitution
	source := extractSourceArray(text, variables)

	return &PKGBUILDInfo{
		Name:    pkgname,
		Version: pkgver,
		URL:     url,
		Source:  source,
	}, nil
}

func parseVariables(text string) Variables {
	variables := make(Variables)

	// Match variable assignments: variable=value
	re := regexp.MustCompile(`^(\w+)\s*=\s*(.+)$`)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			varName := matches[1]
			varValue := cleanValue(matches[2])
			variables[varName] = varValue
		}
	}

	// Now substitute variables within variables (handle nested substitutions)
	for varName, varValue := range variables {
		substitutedValue := substituteVariables(varValue, variables)
		variables[varName] = substitutedValue
	}

	return variables
}

func substituteVariables(value string, variables Variables) string {
	// Replace variables in the format $_variable or ${variable}
	result := value

	// Replace $_variable format
	for varName, varValue := range variables {
		pattern := fmt.Sprintf(`\$%s\b`, varName)
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, varValue)
	}

	// Replace ${variable} format
	for varName, varValue := range variables {
		pattern := fmt.Sprintf(`\$\{%s\}`, varName)
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, varValue)
	}

	return result
}

func extractSourceArray(text string, variables Variables) []string {
	// Find source array definition
	re := regexp.MustCompile(`source\s*=\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil
	}

	sourceContent := matches[1]

	// Split by newlines and clean up
	lines := strings.Split(sourceContent, "\n")
	var sources []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove quotes and clean up
		line = strings.Trim(line, `"'`)
		line = strings.TrimSpace(line)

		if line != "" {
			// Substitute variables
			line = substituteVariables(line, variables)
			sources = append(sources, line)
		}
	}

	return sources
}

func extractValue(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func cleanValue(value string) string {
	// Remove quotes and trim whitespace
	value = strings.Trim(value, `"'`)
	return strings.TrimSpace(value)
}

func getLatestGitHubVersion(repoURL string) (string, error) {
	// Extract owner/repo from GitHub URL
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)
	matches := re.FindStringSubmatch(repoURL)
	if len(matches) < 3 {
		return "", fmt.Errorf("invalid GitHub URL format")
	}

	owner := matches[1]
	repo := matches[2]

	// Try releases first, then tags if releases fails
	version, err := getLatestFromReleases(owner, repo)
	if err != nil {
		// Fallback to tags
		version, err = getLatestFromTags(owner, repo)
		if err != nil {
			return "", fmt.Errorf("failed to fetch version from both releases and tags: %w", err)
		}
	}

	return version, nil
}

func getLatestFromReleases(owner, repo string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from GitHub releases API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub releases API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Remove 'v' prefix if present
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

func getLatestFromTags(owner, repo string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from GitHub tags API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub tags API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var tags []GitHubTag
	if err := json.Unmarshal(body, &tags); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	// Get the first (latest) tag
	version := strings.TrimPrefix(tags[0].Name, "v")
	return version, nil
}

func checkAndSuggestURLFormat(sourceUrls []string, currentUrl string) {
	for _, sourceUrl := range sourceUrls {
		// Check if this is an archive URL (contains archive/refs/tags/)
		if strings.Contains(sourceUrl, "archive/refs/tags/") {
			// Extract the base URL from the current URL
			baseURL := strings.TrimSuffix(currentUrl, "/")

			fmt.Printf("⚠️  Source URL uses archive format: %s\n", sourceUrl)
			fmt.Printf("   Consider using tags endpoint: %s/tags\n", baseURL)
			fmt.Printf("   Example: %s/tags\n", baseURL)
			fmt.Println("")

			// Check tags page for newer versions
			checkTagsPageForNewerVersion(baseURL, sourceUrl)
		}
	}
}

func isPackageIgnored(packageName string) bool {
	ignoreFile := ".aurvtignore"
	
	// Check if .aurvtignore exists
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		return false
	}
	
	// Read .aurvtignore file
	file, err := os.Open(ignoreFile)
	if err != nil {
		return false
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Check if package name matches
		if line == packageName {
			return true
		}
	}
	
	return false
}

func checkTagsPageForNewerVersion(baseURL, sourceUrl string) {
	// Extract current version from source URL
	// Format: filename-version.tar.gz::url/archive/refs/tags/version.tar.gz
	re := regexp.MustCompile(`archive/refs/tags/([^/]+)\.tar\.gz`)
	matches := re.FindStringSubmatch(sourceUrl)
	if len(matches) < 2 {
		fmt.Println("   Could not extract current version from source URL")
		return
	}

	currentVersion := matches[1]

	// Extract owner/repo from base URL
	urlRe := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)
	urlMatches := urlRe.FindStringSubmatch(baseURL)
	if len(urlMatches) < 3 {
		fmt.Println("   Could not extract owner/repo from URL")
		return
	}

	owner := urlMatches[1]
	repo := urlMatches[2]

	// Get latest version from tags API
	latestVersion, err := getLatestFromTags(owner, repo)
	if err != nil {
		// Instead of showing an error, just return silently
		// This handles cases where the tags API might not be available or accessible
		return
	}

	// Compare versions
	if latestVersion != currentVersion {
		fmt.Printf("   🔄 New version available on tags page: %s → %s\n", currentVersion, latestVersion)
		fmt.Printf("   Tags page: %s/tags\n", baseURL)
	} else {
		fmt.Printf("   ✅ Current version %s is up to date on tags page\n", currentVersion)
	}
}
