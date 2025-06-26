package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jparound30/createReviewRecordSheet/pkg/backlog"
)

// CLI はコマンドラインインターフェースを表します
type CLI struct {
	reader *bufio.Reader
	client *backlog.Client
}

// NewCLI は新しいCLIを作成します
func NewCLI(client *backlog.Client) *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		client: client,
	}
}

// SelectProject はユーザーにプロジェクトを選択するよう促します
func (c *CLI) SelectProject() (*backlog.Project, error) {
	// プロジェクトを取得
	projects, err := c.client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("error getting projects: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("no projects found")
	}

	// プロジェクトを表示
	fmt.Println("Available projects:")
	for i, project := range projects {
		fmt.Printf("%d. %s\n", i+1, project.Name)
	}

	// ユーザーにプロジェクトの選択を促す
	fmt.Print("Select a project (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// 入力を解析
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 入力を検証
	if index < 1 || index > len(projects) {
		return nil, fmt.Errorf("invalid project number: %d", index)
	}

	return &projects[index-1], nil
}

// SelectRepository はユーザーにリポジトリを選択するよう促します
func (c *CLI) SelectRepository(projectID int) (*backlog.Repository, error) {
	// リポジトリを取得
	repositories, err := c.client.GetRepositories(projectID)
	if err != nil {
		return nil, fmt.Errorf("error getting repositories: %w", err)
	}

	if len(repositories) == 0 {
		return nil, fmt.Errorf("no repositories found for the selected project")
	}

	// リポジトリを表示
	fmt.Println("Available repositories:")
	for i, repo := range repositories {
		fmt.Printf("%d. %s\n", i+1, repo.Name)
	}

	// ユーザーにリポジトリの選択を促す
	fmt.Print("Select a repository (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// 入力を解析
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 入力を検証
	if index < 1 || index > len(repositories) {
		return nil, fmt.Errorf("invalid repository number: %d", index)
	}

	return &repositories[index-1], nil
}

// SelectPullRequest はユーザーにプルリクエストを選択するよう促します
func (c *CLI) SelectPullRequest(projectID int, repoID int) (*backlog.PullRequest, error) {
	// プルリクエストを取得
	pullRequests, err := c.client.GetPullRequests(projectID, repoID)
	if err != nil {
		return nil, fmt.Errorf("error getting pull requests: %w", err)
	}

	if len(pullRequests) == 0 {
		return nil, fmt.Errorf("no pull requests found for the selected repository")
	}

	// プルリクエストを表示
	fmt.Println("Available pull requests:")
	for i, pr := range pullRequests {
		fmt.Printf("%d. %s\n", i+1, pr.Summary)
	}

	// ユーザーにプルリクエストの選択を促す
	fmt.Print("Select a pull request (enter number): ")
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// 入力を解析
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 入力を検証
	if index < 1 || index > len(pullRequests) {
		return nil, fmt.Errorf("invalid pull request number: %d", index)
	}

	return &pullRequests[index-1], nil
}
