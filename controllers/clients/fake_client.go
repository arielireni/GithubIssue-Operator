package clients

import (
	"net/http"
	"os"
)

///////////////////////////////////////////////////////////////////////
/////////////// Implementation of GitHubClient - "test" ///////////////
///////////////////////////////////////////////////////////////////////

type FakeClient struct {
	HttpClient http.Client
	Token      string
	RepoURL    string
}

func (g *FakeClient) InitDataStructs(repo, title, body string) (*Repo, *Issue, *Details) {
	return nil, nil, nil
}

func (g *FakeClient) FindIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error) {
	return nil, nil
}

func (g *FakeClient) CreateNewIssue(issueData *Issue, detailsData *Details) (*Issue, *Error) {
	return nil, nil
}

func (g *FakeClient) EditIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error {
	return nil
}

func (g *FakeClient) CloseIssue(issue *Issue, issueData *Issue, detailsData *Details) *Error {
	return nil
}

func NewFakeClient(repoURL string) *FakeClient {
	return &FakeClient{
		HttpClient: http.Client{},
		Token:      os.Getenv("TOKEN"),
		RepoURL:    repoURL,
	}
}
