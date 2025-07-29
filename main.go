package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type PKGBUILDInfo struct {
	Name    string
	Version string
	URL     string
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aur-version-tool [package-directory]",
		Short: "Check for newer versions of AUR packages on GitHub",
		Long:  `A CLI tool to check if newer versions are available for AUR packages hosted on GitHub.`,
		Args:  cobra.ExactArgs(1),
		Run:   runVersionCheck,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runVersionCheck(cmd *cobra.Command, args []string) {
	packageDir := args[0]
	pkgbuildPath := fmt.Sprintf("%s/PKGBUILD", packageDir)

	// Parse PKGBUILD
	pkgInfo, err := parsePKGBUILD(pkgbuildPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PKGBUILD: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Package: %s\n", pkgInfo.Name)
	fmt.Printf("Current version: %s\n", pkgInfo.Version)
	fmt.Printf("Repository URL: %s\n", pkgInfo.URL)

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
	
	// Extract pkgname
	pkgname := extractValue(text, `pkgname\s*=\s*(.+)`)
	if pkgname == "" {
		return nil, fmt.Errorf("could not find pkgname in PKGBUILD")
	}

	// Extract pkgver
	pkgver := extractValue(text, `pkgver\s*=\s*(.+)`)
	if pkgver == "" {
		return nil, fmt.Errorf("could not find pkgver in PKGBUILD")
	}

	// Extract url
	url := extractValue(text, `url\s*=\s*(.+)`)
	if url == "" {
		return nil, fmt.Errorf("could not find url in PKGBUILD")
	}

	return &PKGBUILDInfo{
		Name:    cleanValue(pkgname),
		Version: cleanValue(pkgver),
		URL:     cleanValue(url),
	}, nil
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
	
	// Make API request
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
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