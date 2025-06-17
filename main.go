package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jparound30/createReviewRecordSheet/internal/config"
	"github.com/jparound30/createReviewRecordSheet/pkg/backlog"
	"github.com/jparound30/createReviewRecordSheet/pkg/cli"
	"github.com/jparound30/createReviewRecordSheet/pkg/excel"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create Backlog client
	client := backlog.NewClient(cfg)

	// Create CLI
	cliApp := cli.NewCLI(client)

	// Create Excel generator
	excelGen := excel.NewGenerator()

	// Select project
	fmt.Println("Fetching projects...")
	project, err := cliApp.SelectProject()
	if err != nil {
		log.Fatalf("Error selecting project: %v", err)
	}
	fmt.Printf("Selected project: %s\n", project.Name)

	// Select repository
	fmt.Println("Fetching repositories...")
	repo, err := cliApp.SelectRepository(project.ID)
	if err != nil {
		log.Fatalf("Error selecting repository: %v", err)
	}
	fmt.Printf("Selected repository: %s\n", repo.Name)

	// Select pull request
	fmt.Println("Fetching pull requests...")
	pr, err := cliApp.SelectPullRequest(project.ID, repo.ID)
	if err != nil {
		log.Fatalf("Error selecting pull request: %v", err)
	}
	fmt.Printf("Selected pull request: %s\n", pr.Summary)

	// Get comments
	fmt.Println("Fetching comments...")
	comments, err := client.GetComments(project.ID, repo.ID, pr.Number)
	if err != nil {
		log.Fatalf("Error getting comments: %v", err)
	}
	fmt.Printf("Found %d comments\n", len(comments))

	// Generate Excel file
	fmt.Println("Generating Excel file...")
	filename, err := excelGen.GenerateCommentSheet(
		project.Name,
		repo.Name,
		pr.Summary,
		comments,
		client,
	)
	if err != nil {
		log.Fatalf("Error generating Excel file: %v", err)
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	fmt.Printf("Excel file generated successfully: %s\\%s\n", currentDir, filename)
}
