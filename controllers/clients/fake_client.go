package clients

import (
	"fmt"
	"strings"
)

/* Implementation of FakeClient - "test" */

type FakeClient struct {
	issues []Issue
	err    error
}

func (f *FakeClient) InitDataStructs(repo, title, body string) (*Repo, *Issue, *Details) {
	splitRepo := strings.Split(repo, "/")
	// Init Repo data
	owner := splitRepo[1]
	repoName := splitRepo[0]
	repoData := Repo{Owner: owner, Repo: repoName}
	// Init Issue data
	issueData := Issue{Title: title, Description: body}
	// Init Details data
	apiURL := "TestURL"
	token := "TestToken"
	detailsData := Details{ApiURL: apiURL, Token: token}
	return &repoData, &issueData, &detailsData
}

func (f *FakeClient) FindIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error) {
	returnErr := Error{}
	for _, issue := range f.issues {
		if issue.Title == issueData.Title {
			return &issue, &returnErr
		}
	}
	returnErr.ErrorCode = fmt.Errorf("FindIssue error")
	returnErr.Message = "Error with find issue"
	return nil, &returnErr
}

func (f *FakeClient) CreateIssue(issueData *Issue, detailsData *Details) (*Issue, *Error) {
	returnErr := Error{
		ErrorCode: f.err,
		Message:   "Error with create issue",
	}
	if f.err != nil {
		return &Issue{}, &returnErr
	}
	newIssue := Issue{
		Title:               issueData.Title,
		Description:         issueData.Description,
		Number:              issueData.Number,
		State:               issueData.State,
		LastUpdateTimestamp: issueData.LastUpdateTimestamp,
	}
	f.issues = append(f.issues, newIssue)
	return &newIssue, nil
}

func (f *FakeClient) EditIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error {
	returnErr := Error{}
	for _, issue := range f.issues {
		if issue.Title == issueData.Title {
			issue.Description = issueData.Description
			return &returnErr
		}
	}
	returnErr.ErrorCode = fmt.Errorf("EditIssue error")
	returnErr.Message = "Error with edit issue"
	return &returnErr
}

func (f *FakeClient) CloseIssue(issue *Issue, issueData *Issue, detailsData *Details) *Error {
	returnErr := Error{}
	for _, issue := range f.issues {
		if issue.Title == issueData.Title {
			issue.State = "closed"
			return &returnErr
		}
	}
	returnErr.ErrorCode = fmt.Errorf("CloseIssue error")
	returnErr.Message = "Error with close issue"
	return &returnErr
}

func NewFakeClient(issues []Issue, isSuccessful bool, err error) *FakeClient {
	if isSuccessful {
		return &FakeClient{
			issues: issues,
		}
	}
	return &FakeClient{
		issues: issues,
		err:    err,
	}
}
