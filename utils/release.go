/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"github.com/Masterminds/semver"
)

// Release represents a GitHub API release
type Release struct {
	ID 			int 						`json:"id"`
	Version *semver.Version `json:"name"`
	Date 		*time.Time 			`json:"created_at"`
	Body 		string 					`json:"body"`
}

var remoteMatcher, _ = regexp.Compile("(?i)^.+github\\.com[:/](.+)\\.git$")

// GitHubReleaseAPICall performs a GitHub release API call.
func GitHubReleaseAPICall(message, method, path, token string, data map[string]string, allowErrors bool) *gentleman.Response {
	cli := gentleman.New()
  cli.URL("https://api.github.com")

	req := cli.Request()
	req.Method(method)
	req.Path(path)
	
	if token != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if method != "GET" {
		req.Use(body.JSON(data))
	}
	
  // Perform the request
  res, err := req.Send()

	if err != nil {
		Fatal("Cannot %s due to a network error: {errorPrimary}%s{-}", message, err.Error())
  }

	if !res.Ok {
		if res.StatusCode == 401 { 
			Fatal("Cannot %s due to an authentication error.", message)
		} else if !allowErrors {
			Fatal(
				"Cannot %s due to an HTTP error: {secondary}[HTTP %d]{-} {errorPrimary}%s{-}", 
				message, res.StatusCode, res.String(),
			)
		}
	}

	return res
}

// DetectGithubRepository detects the GitHub repository.
func DetectGithubRepository(remote string, allowFailure bool) string {
	result := Execute(false, "git", "remote", "get-url", remote)
	result.Verify("git", "Cannot get GIT remote url")

	remoteURL := strings.TrimSpace(result.Stdout)
	if !strings.HasPrefix(remoteURL, "https://github.com") && !strings.HasPrefix(remoteURL, "git@github.com") {
		if allowFailure {
			return ""
		}

		Fatal("The GIT remote {errorPrimary}%s{-} is not a GitHub repository.", remote)
	}

	return remoteMatcher.FindStringSubmatch(remoteURL)[1]
}

// FindRelease finds a release id on GitHub API.
func FindRelease(repository, token, version string) int {
	res := GitHubReleaseAPICall(
		"find a GitHub release", "GET", fmt.Sprintf("/repos/%s/releases/tags/v%s", repository, version), 
		token, map[string]string{}, true,
	)

	if res.StatusCode == 404 {
		return 0
	} 

	var release Release
	err := res.JSON(&release)
	
	if err != nil {
		Fatal("Cannot decode JSON response to %s: {errorPrimary}%s{-}", err.Error())
	}	

	return release.ID
}

// SaveRelease creates or updates a release on GitHub 
func SaveRelease(version *semver.Version, repository, remote, token string, dryRun bool) {
	// Get and format changes
	versions := GetVersions()

	currentIndex := -1
	for i, v := range versions {
		if v.Equal(version) {
			currentIndex = i
			break
		}
	}

	var changes []Change

	if currentIndex > 0 {
		previousVersion := versions[currentIndex-1]
		changes = ListChanges(version.String(), previousVersion.String())
	} else {
		changes = ListChanges(version.String(), GetFirstCommitHash())
	}

	changelog := strings.TrimSpace(FormatReleaseChanges(repository, changes))
	data := map[string]string{
		"tag_name": fmt.Sprintf("v%s", version.String()),
		"name": version.String(),
		"body": changelog,
	}

	// Check if a release exists
	existing := FindRelease(repository, "", version.String())

	// Perform the right operation on GitHub
	if existing != 0 {
		if NotifyStep(dryRun, "", "Will update", "Updating", " GitHub release {primary}%s{-}...", version.String()) {
			GitHubReleaseAPICall(
				"update a GitHub release", "PATCH", fmt.Sprintf("/repos/%s/releases/%d", repository, existing), 
				token, data, false,
			)
		}
	} else {
		if NotifyStep(dryRun, "", "Will create", "Creating", " GitHub release {primary}%s{-}...", version.String()) {
			GitHubReleaseAPICall(
				"create a GitHub release", "POST", fmt.Sprintf("/repos/%s/releases", repository), 
				token, data, false,
			)
		}
	}
}