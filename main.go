package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	http.HandleFunc("/", requestHandler)

	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appVersion := r.URL.Query().Get("app_version")
	location := r.URL.Query().Get("location")

	if appVersion == "" || location == "" {
		http.Error(w, "Missing required parameters: app_version and location", http.StatusBadRequest)
		return
	}

	isDirty := strings.Contains(appVersion, "-dirty-")
	gitRef := extractGitRef(appVersion)
	if gitRef == "" {
		http.Error(w, "Could not extract git reference from app_version", http.StatusBadRequest)
		return
	}

	repoName := extractRepoName(location, appVersion)
	if repoName == "" {
		http.Error(w, "Could not determine repository name from location and appVersion", http.StatusBadRequest)
		return
	}

	filePath, lineNum := extractLocation(location)
	if filePath == "" {
		http.Error(w, "Could not extract file path from location", http.StatusBadRequest)
		return
	}

	githubURL := fmt.Sprintf("https://github.com/VictoriaMetrics/%s/blob/%s/%s", repoName, gitRef, filePath)
	if lineNum != "" {
		githubURL += fmt.Sprintf("#L%s", lineNum)
	}

	if isDirty {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>Dirty Build Warning</title></head>
<body>
<p>Warning: %s is a dirty build. The source code may not exactly match the running code.</p>
<p><a href="%s">Continue to GitHub</a></p>
</body>
</html>`, appVersion, githubURL)

		fmt.Fprint(w, html)
		return
	}

	log.Printf("Redirecting to: %s", githubURL)
	http.Redirect(w, r, githubURL, http.StatusFound)
}

func extractRepoName(location, appVersion string) string {
	repo, _, found := strings.Cut(location, "/")
	if !found {
		return ""
	}

	switch repo {
	case "VictoriaMetrics", "VictoriaLogs", "VictoriaTraces":
		// valid base repo names
	default:
		return ""
	}

	if strings.Contains(appVersion, "enterprise") {
		repo += "-enterprise"
	}

	return repo
}

func extractGitRef(appVersion string) string {
	// Check for tags (e.g., tags-v1.129.1-0-g5e98e0cff5)
	if strings.Contains(appVersion, "-tags-") {
		re := regexp.MustCompile(`-tags-(.*?)-\d+-g[0-9a-f]+?`)
		matches := re.FindStringSubmatch(appVersion)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// Check for branch at head with 0 commits ahead (e.g., heads-master-0-g1db7597e45)
	if strings.Contains(appVersion, "-heads-") {
		re := regexp.MustCompile(`-heads-(.*?)-\d+-g[0-9a-f]+?`)
		matches := re.FindStringSubmatch(appVersion)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// Default: extract commit hash
	re := regexp.MustCompile(`-g([0-9a-f]+)`)
	matches := re.FindStringSubmatch(appVersion)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func extractLocation(location string) (filePath, lineNum string) {
	// Format: VictoriaMetrics/lib/vmselectapi/server.go:200
	location = strings.TrimPrefix(location, "VictoriaMetrics/")

	parts := strings.Split(location, ":")
	if len(parts) == 0 {
		return "", ""
	}

	filePath = parts[0]
	if len(parts) > 1 {
		lineNum = parts[1]
	}

	return filePath, lineNum
}
