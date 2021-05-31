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

///////////////////////////////////////////////////////////////////////
//////////// Implementation of GitHubClient - "production" ////////////
///////////////////////////////////////////////////////////////////////

type GithubClient struct {
	HttpClient http.Client
	Token      string
	RepoURL    string
}

// InitDataStructs initializes issueData & detailsData
func (g *GithubClient) InitDataStructs(repo, title, body string) (*Repo, *Issue, *Details) {
	splittedRepo := strings.Split(repo, "/")
	owner := splittedRepo[1]
	repoName := splittedRepo[0]
	repoData := Repo{Owner: owner, Repo: repoName}

	issueData := Issue{Title: title, Description: body}

	apiURL := "https://api.clients.com/repos/" + repo + "/issues"
	token := os.Getenv("TOKEN")
	detailsData := Details{ApiURL: apiURL, Token: token}

	return &repoData, &issueData, &detailsData
}

func (g *GithubClient) FindIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error) {
	apiURL := detailsData.ApiURL + "?state=all"
	/* API request for all repository's issues */
	jsonData, _ := json.Marshal(&repoData)
	// creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		typedErr := Error{ErrorCode: err, Message: "GET request from GitHub API faild"}
		return nil, &typedErr
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// print body as it may contain hints in case of errors
	fmt.Println(string(body))

	// Create array with all repository's issues
	var allIssues []Issue
	returnErr := Error{}
	err = json.Unmarshal(body, &allIssues)
	if err != nil {
		typedErr := Error{ErrorCode: err, Message: "Unmarshal faild"}
		return nil, &typedErr
	}

	// If we found the issue, we'll return it. Otherwise, return nil
	for _, issue := range allIssues {
		if issue.Title == issueData.Title {
			return &issue, nil
		}
	}
	return nil, nil
}

func (g *GithubClient) CreateIssue(issueData *Issue, detailsData *Details) (*Issue, *Error) {
	apiURL := detailsData.ApiURL
	// make it json
	jsonData, _ := json.Marshal(issueData)
	// creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Printf("fatal error")
		//log.Fatal(err)
		returnErr := Error{ErrorCode: err, Message: "POST request from GitHub API faild"}
		return nil, &returnErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		//log.Fatal(err)
		returnErr := Error{ErrorCode: err, Message: "Creating GitHub issue faild"}
		return nil, &returnErr
	}
	var issue *Issue
	issueBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(issueBody, &issue)
	return issue, nil
}

func (g *GithubClient) EditIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error {
	issue.Description = issueData.Description
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	fmt.Printf("URL: " + issueApiURL)
	jsonData, _ := json.Marshal(issue)

	// Now update
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		return &Error{ErrorCode: err, Message: "PATCH request from GitHub API faild"}
		//log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		//log.Fatal(err)
		return &Error{ErrorCode: err, Message: "Editing GitHub issue faild"}
	}
	return nil
}

func (g *GithubClient) CloseIssue(issue *Issue, issueData *Issue, detailsData *Details) *Error {
	issueApiURL := detailsData.ApiURL + "/" + fmt.Sprint(issue.Number)
	issue.State = "closed"
	jsonData, _ := json.Marshal(issue)

	// Now update
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", issueApiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+detailsData.Token)
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(err)
		return &Error{ErrorCode: err, Message: "PATCH request from GitHub API faild"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		// print body as it may contain hints in case of errors
		fmt.Println(string(body))
		//log.Fatal(err)
		return &Error{ErrorCode: err, Message: "Closing GitHub issue faild"}
	}
	return nil
}

func NewGithubClient(repoURL string) *GithubClient {
	return &GithubClient{
		HttpClient: http.Client{},
		Token:      os.Getenv("TOKEN"),
		RepoURL:    repoURL,
	}
}
