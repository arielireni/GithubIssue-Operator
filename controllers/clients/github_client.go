package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/* Implementation of GitHubClient - "production" */

type GithubClient struct {
	HttpClient http.Client
	Token      string
	RepoURL    string
}

// InitDataStructs initializes issueData & detailsData
func (g *GithubClient) InitDataStructs(repo, title, body string) (*Repo, *Issue, *Details) {
	splitRepo := strings.Split(repo, "/")
	// Init Repo data
	owner := splitRepo[1]
	repoName := splitRepo[0]
	repoData := Repo{Owner: owner, Repo: repoName}
	// Init Issue data
	issueData := Issue{Title: title, Description: body}
	// Init Details data
	apiURL := "https://api.github.com/repos/" + repo + "/issues"
	token := os.Getenv("TOKEN")
	detailsData := Details{ApiURL: apiURL, Token: token}

	return &repoData, &issueData, &detailsData
}

func (g *GithubClient) FindIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error) {
	apiURL := detailsData.ApiURL + "?state=all"
	// API request for all repository's issues
	jsonData, _ := json.Marshal(&repoData)
	// Creating client to set custom headers for Authorization
	client := g.HttpClient
	req, _ := http.NewRequest("GET", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	returnErr := Error{}
	if err != nil {
		returnErr = Error{ErrorCode: err, Message: "GET request from GitHub API failed"}
		return nil, &returnErr
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// Create array with all repository's issues
	var allIssues []Issue
	err = json.Unmarshal(body, &allIssues)
	if err != nil {
		returnErr = Error{ErrorCode: err, Message: "Unmarshal failed"}
		return nil, &returnErr
	}
	// If we found the issue, we will return it. Otherwise, return nil
	for _, issue := range allIssues {
		if issue.Title == issueData.Title {
			return &issue, &returnErr
		}
	}
	return nil, &returnErr
}

func (g *GithubClient) CreateIssue(issueData *Issue, detailsData *Details) (*Issue, *Error) {
	apiURL := detailsData.ApiURL
	// Make it json
	jsonData, _ := json.Marshal(issueData)
	// Creating client to set custom headers for Authorization
	client := g.HttpClient
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	returnErr := Error{}
	if err != nil {
		returnErr = Error{ErrorCode: err, Message: "POST request from GitHub API failed"}
		return nil, &returnErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		returnErr = Error{ErrorCode: err, Message: "Creating GitHub issue failed"}
		return nil, &returnErr
	}
	var issue *Issue
	issueBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(issueBody, &issue)
	return issue, &returnErr
}

func (g *GithubClient) EditIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error {
	issue.Description = issueData.Description
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	jsonData, _ := json.Marshal(issue)
	// Now update
	client := g.HttpClient
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	returnErr := Error{}
	if err != nil {
		returnErr = Error{ErrorCode: err, Message: "PATCH request from GitHub API failed"}
		return &returnErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		returnErr = Error{ErrorCode: err, Message: "Editing GitHub issue failed"}
		return &returnErr
	}
	return &returnErr
}

func (g *GithubClient) CloseIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error {
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	issue.State = "closed"
	jsonData, _ := json.Marshal(issue)
	// Now update
	client := g.HttpClient
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	returnErr := Error{}
	if err != nil {
		returnErr = Error{ErrorCode: err, Message: "PATCH request from GitHub API failed"}
		return &returnErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		returnErr = Error{ErrorCode: err, Message: "Closing GitHub issue failed"}
		return &returnErr
	}
	return &returnErr
}

func NewGithubClient(repoURL string) GithubClient {
	return GithubClient{
		HttpClient: http.Client{},
		Token:      os.Getenv("TOKEN"),
		RepoURL:    repoURL,
	}
}
