package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func (a *Adapter) GetRepoURL() string {
	remote, err := a.repo.Remote("origin")
	if err != nil {
		return ""
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return ""
	}

	return normalizeGitURL(urls[0])
}

func (a *Adapter) GetRepoName() string {
	urlStr := a.GetRepoURL()
	if urlStr == "" {
		return "unknown"
	}

	parts := strings.Split(urlStr, "/")
	if len(parts) > 0 {
		return strings.TrimSuffix(parts[len(parts)-1], ".git")
	}
	return "unknown"
}

func normalizeGitURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "http") {
		return strings.TrimSuffix(rawURL, ".git")
	}

	// Handle SSH URLs (git@github.com:user/repo.git)
	re := regexp.MustCompile(`^git@([^:]+):(.+)$`)
	matches := re.FindStringSubmatch(rawURL)
	if len(matches) == 3 {
		host := matches[1]
		path := matches[2]
		return fmt.Sprintf("https://%s/%s", host, strings.TrimSuffix(path, ".git"))
	}

	// Handle SCP-like syntax without user
	if strings.Contains(rawURL, ":") && !strings.Contains(rawURL, "://") {
		parts := strings.SplitN(rawURL, ":", 2)
		return fmt.Sprintf("https://%s/%s", parts[0], strings.TrimSuffix(parts[1], ".git"))
	}

	return rawURL
}

func (a *Adapter) buildCommitURL(hash string) string {
	baseURL := a.GetRepoURL()
	if baseURL == "" {
		return ""
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// GitHub/GitLab style
	if strings.Contains(u.Host, "github.com") || strings.Contains(u.Host, "gitlab.com") {
		return fmt.Sprintf("%s/commit/%s", baseURL, hash)
	}

	// Bitbucket style
	if strings.Contains(u.Host, "bitbucket.org") {
		return fmt.Sprintf("%s/commits/%s", baseURL, hash)
	}

	return fmt.Sprintf("%s/commit/%s", baseURL, hash)
}

func (a *Adapter) GetDiffURL(base, target string) string {
	baseURL := a.GetRepoURL()
	if baseURL == "" {
		return ""
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	if strings.Contains(u.Host, "github.com") || strings.Contains(u.Host, "gitlab.com") {
		return fmt.Sprintf("%s/compare/%s...%s", baseURL, base[:7], target[:7])
	}

	if strings.Contains(u.Host, "bitbucket.org") {
		return fmt.Sprintf("%s/branches/compare/%s..%s", baseURL, target[:7], base[:7])
	}

	return fmt.Sprintf("%s/compare/%s...%s", baseURL, base[:7], target[:7])
}
