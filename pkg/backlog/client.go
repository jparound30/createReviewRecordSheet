package backlog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jparound30/createReviewRecordSheet/internal/config"
)

// Client represents a Backlog API client
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// NewClient creates a new Backlog API client
func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Project represents a Backlog project
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Repository represents a Backlog Git repository
type Repository struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PullRequest represents a Backlog pull request
type PullRequest struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	Summary string `json:"summary"`
}

// Comment represents a Backlog comment
type Comment struct {
	ID            int           `json:"id"`
	OldBlobID     *string       `json:"oldBlobId"`
	NewBlobID     *string       `json:"newBlobId"`
	FilePath      *string       `json:"filePath"`
	Position      *int          `json:"position"`
	Content       string        `json:"content"`
	ChangeLog     []ChangeLog   `json:"changeLog"`
	CreatedUser   User          `json:"createdUser"`
	Created       time.Time     `json:"created"`
	Updated       time.Time     `json:"updated"`
	Stars         []interface{} `json:"stars"`
	Notifications []interface{} `json:"notifications"`
}

// NulabAccount represents Nulab account information
type NulabAccount struct {
	NulabID  string `json:"nulabId"`
	Name     string `json:"name"`
	UniqueID string `json:"uniqueId"`
	IconURL  string `json:"iconUrl"`
}

// User represents a Backlog user
type User struct {
	ID            int          `json:"id"`
	UserID        string       `json:"userId"`
	Name          string       `json:"name"`
	RoleType      int          `json:"roleType"`
	Lang          string       `json:"lang"`
	MailAddress   string       `json:"mailAddress"`
	NulabAccount  NulabAccount `json:"nulabAccount"`
	Keyword       string       `json:"keyword"`
	LastLoginTime string       `json:"lastLoginTime"`
}

// ChangeLog represents a change in a comment
type ChangeLog struct {
	Field          string `json:"field"`
	OriginalValue  string `json:"originalValue"`
	NewValue       string `json:"newValue"`
	AttachmentInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"attachmentInfo"`
	AttributeInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"attributeInfo"`
	NotificationInfo struct {
		Type string `json:"type"`
	} `json:"notificationInfo"`
}

// CommentLocation represents the location of a comment
type CommentLocation struct {
	Path     string
	LineFrom int
	LineTo   int
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() ([]Project, error) {
	endpoint := fmt.Sprintf("%s/projects?apiKey=%s", c.config.GetBacklogBaseURL(), c.config.BacklogAPIKey)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting projects: %s, status code: %d", string(body), resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("error decoding projects: %w", err)
	}

	return projects, nil
}

// GetRepositories retrieves all repositories for a project
func (c *Client) GetRepositories(projectID int) ([]Repository, error) {
	endpoint := fmt.Sprintf("%s/projects/%d/git/repositories?apiKey=%s",
		c.config.GetBacklogBaseURL(), projectID, c.config.BacklogAPIKey)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting repositories: %s, status code: %d", string(body), resp.StatusCode)
	}

	var repositories []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		return nil, fmt.Errorf("error decoding repositories: %w", err)
	}

	return repositories, nil
}

// GetPullRequests retrieves all pull requests for a repository
func (c *Client) GetPullRequests(projectID int, repoID int) ([]PullRequest, error) {
	endpoint := fmt.Sprintf("%s/projects/%d/git/repositories/%d/pullRequests?apiKey=%s&statusId[]=1&statusId[]=2&statusId[]=3",
		c.config.GetBacklogBaseURL(), projectID, repoID, c.config.BacklogAPIKey)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting pull requests: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting pull requests: %s, status code: %d", string(body), resp.StatusCode)
	}

	var pullRequests []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, fmt.Errorf("error decoding pull requests: %w", err)
	}

	return pullRequests, nil
}

// GetComments retrieves all comments for a pull request
func (c *Client) GetComments(projectID int, repoID int, pullRequestID int) ([]Comment, error) {
	endpoint := fmt.Sprintf("%s/projects/%d/git/repositories/%d/pullRequests/%d/comments?apiKey=%s&count=100",
		c.config.GetBacklogBaseURL(), projectID, repoID, pullRequestID, c.config.BacklogAPIKey)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting comments: %s, status code: %d", string(body), resp.StatusCode)
	}

	var comments []Comment
	if err := json.NewDecoder(strings.NewReader(string(body))).Decode(&comments); err != nil {
		return nil, fmt.Errorf("error decoding comments: %w", err)
	}

	return comments, nil
}

// GetCommentLocation extracts the location information from a comment
func (c *Client) GetCommentLocation(comment Comment) string {
	if comment.FilePath == nil || comment.Position == nil {
		return "全体"
	}
	return fmt.Sprintf("%s: %d行目", *comment.FilePath, *comment.Position)
}

// FormatDate formats a time.Time as yyyy/mm/dd
func FormatDate(t time.Time) string {
	return t.Format("2006/01/02")
}
