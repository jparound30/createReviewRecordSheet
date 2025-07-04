package excel

import (
	"fmt"
	"time"

	"github.com/jparound30/createReviewRecordSheet/pkg/backlog"
	"github.com/xuri/excelize/v2"
)

// Generator represents an Excel file generator
type Generator struct{}

// NewGenerator creates a new Excel generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateCommentSheet generates an Excel file with comments
func (g *Generator) GenerateCommentSheet(
	projectName string,
	repoName string,
	pullRequestName string,
	comments []backlog.Comment,
	client *backlog.Client,
) (string, error) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	sheetName := "Comments"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := []string{"コメントがつけられた場所", "ユーザー名", "日付", "コメント内容"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 100)

	// Create a style with top alignment and text wrapping
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating cell style: %w", err)
	}

	// Apply the style to all cells in the sheet
	err = f.SetCellStyle(sheetName, "A1", fmt.Sprintf("D%d", len(comments)+1), style)
	if err != nil {
		return "", fmt.Errorf("error applying cell style: %w", err)
	}

	// Add comments
	for i, comment := range comments {
		row := i + 2 // Start from row 2 (after headers)

		// Get comment location
		location := client.GetCommentLocation(comment)

		// Set values
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), location)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), comment.CreatedUser.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), backlog.FormatDate(comment.Created))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), comment.Content)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s_%s_%s_%s.xlsx",
		projectName,
		repoName,
		pullRequestName,
		now.Format("20060102"))

	// Save the file
	if err := f.SaveAs(filename); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	return filename, nil
}
