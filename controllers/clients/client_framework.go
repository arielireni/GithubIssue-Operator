package clients

type ClientFrame interface {
	InitDataStructs(repo, title, body string) (*Repo, *Issue, *Details)
	FindIssue(repoData *Repo, issueData *Issue, detailsData *Details) (*Issue, *Error)
	CreateIssue(issueData *Issue, detailsData *Details) (*Issue, *Error)
	EditIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error
	CloseIssue(issueData *Issue, issue *Issue, detailsData *Details) *Error
}

// Repo structure declaration - all data fields for getting a repo's issues list
type Repo struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

// Issue structure declaration - all data fields for a new clients issue submission
type Issue struct {
	Title               string `json:"title"`
	Description         string `json:"body"`
	Number              int    `json:"number"`
	State               string `json:"state,omitempty"`
	LastUpdateTimestamp string `json:"updated_at,omitempty"`
}

// Details structure declaration - all owner's details
type Details struct {
	ApiURL string
	Token  string
}

// Error structure declaration - errors with messages
type Error struct {
	ErrorCode error
	Message   string
}
