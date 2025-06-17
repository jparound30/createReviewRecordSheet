package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jparound30/createReviewRecordSheet/pkg/backlog"
)

// CLI represents the command-line interface
type CLI struct {
	reader *bufio.Reader
	client *backlog.Client
}

// NewCLI creates a new CLI
func NewCLI(client *backlog.Client) *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		client: client,
	}
}

// SelectProject prompts the user to select a project
func (c *CLI) SelectProject() (*backlog.Project, error) {
	// Get projects
	projects, err := c.client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("error getting projects: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("no projects found")
	}

	// Display projects
	fmt.Println("Available projects:")
	for i, project := range projects {
		fmt.Printf("%d. %s\n", i+1, project.Name)
	}

	// Prompt user to select a project
	fmt.Print("Select a project (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// Parse input
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Validate input
	if index < 1 || index > len(projects) {
		return nil, fmt.Errorf("invalid project number: %d", index)
	}

	return &projects[index-1], nil
}

// SelectRepository prompts the user to select a repository
func (c *CLI) SelectRepository(projectID int) (*backlog.Repository, error) {
	// Get repositories
	repositories, err := c.client.GetRepositories(projectID)
	if err != nil {
		return nil, fmt.Errorf("error getting repositories: %w", err)
	}

	if len(repositories) == 0 {
		return nil, fmt.Errorf("no repositories found for the selected project")
	}

	// Display repositories
	fmt.Println("Available repositories:")
	for i, repo := range repositories {
		fmt.Printf("%d. %s\n", i+1, repo.Name)
	}

	// Prompt user to select a repository
	fmt.Print("Select a repository (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// Parse input
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Validate input
	if index < 1 || index > len(repositories) {
		return nil, fmt.Errorf("invalid repository number: %d", index)
	}

	return &repositories[index-1], nil
}

// SelectPullRequest prompts the user to select a pull request
func (c *CLI) SelectPullRequest(projectID int, repoID int) (*backlog.PullRequest, error) {
	// Get pull requests
	pullRequests, err := c.client.GetPullRequests(projectID, repoID)
	if err != nil {
		return nil, fmt.Errorf("error getting pull requests: %w", err)
	}

	if len(pullRequests) == 0 {
		return nil, fmt.Errorf("no pull requests found for the selected repository")
	}

	// Display pull requests
	fmt.Println("Available pull requests:")
	for i, pr := range pullRequests {
		fmt.Printf("%d. %s\n", i+1, pr.Summary)
	}

	// Prompt user to select a pull request
	fmt.Print("Select a pull request (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// Parse input
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Validate input
	if index < 1 || index > len(pullRequests) {
		return nil, fmt.Errorf("invalid pull request number: %d", index)
	}

	return &pullRequests[index-1], nil
}
