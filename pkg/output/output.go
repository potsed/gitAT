package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	// Logger instance for the application
	Logger *log.Logger

	// Styles for different output types
	styles = struct {
		Success  lipgloss.Style
		Error    lipgloss.Style
		Warning  lipgloss.Style
		Info     lipgloss.Style
		Title    lipgloss.Style
		Subtitle lipgloss.Style
		Code     lipgloss.Style
		Dim      lipgloss.Style
	}{
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffaa00")).
			Bold(true),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00aaff")).
			Bold(true),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cccccc")).
			Bold(true),
		Code: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffaa")).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1),
		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
)

// Init initializes the output package with a logger
func Init() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		Level:           log.InfoLevel,
		Prefix:          "GitAT 🚀",
		TimeFormat:      "15:04:05",
	})
}

// Success logs a success message using the logger
func Success(format string, args ...interface{}) {
	Logger.Info(fmt.Sprintf(format, args...))
}

// SuccessWithFields logs a success message with additional fields
func SuccessWithFields(message string, fields map[string]interface{}) {
	// Convert map to key-value pairs for proper logging
	var args []interface{}
	for key, value := range fields {
		args = append(args, key, value)
	}
	Logger.Info(message, args...)
}

// SaveSuccess logs a save operation with beautiful formatting
func SaveSuccess(branch, message string) {
	Success("Changes saved successfully")
	Info("Branch: %s", branch)
	Info("Commit: %s", message)
}

// Error logs an error message using the logger
func Error(format string, args ...interface{}) {
	Logger.Error(fmt.Sprintf(format, args...))
}

// Warning logs a warning message using the logger
func Warning(format string, args ...interface{}) {
	Logger.Warn(fmt.Sprintf(format, args...))
}

// Info logs an info message using the logger
func Info(format string, args ...interface{}) {
	Logger.Info(fmt.Sprintf(format, args...))
}

// Title prints a title (keeping styled output for titles)
func Title(text string) {
	fmt.Println(styles.Title.Render(text))
}

// Subtitle prints a subtitle (keeping styled output for subtitles)
func Subtitle(text string) {
	fmt.Println(styles.Subtitle.Render(text))
}

// Code prints code with syntax highlighting (keeping styled output for code)
func Code(text string) {
	fmt.Println(styles.Code.Render(text))
}

// Dim prints dimmed text (keeping styled output for dimmed text)
func Dim(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(styles.Dim.Render(msg))
}

// Markdown renders markdown text with glamour
func Markdown(text string) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return err
	}

	rendered, err := renderer.Render(text)
	if err != nil {
		return err
	}

	fmt.Print(rendered)
	return nil
}

// Table creates a properly aligned table with lipgloss
func Table(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))

	// Check header widths
	for i, header := range headers {
		if len(header) > colWidths[i] {
			colWidths[i] = len(header)
		}
	}

	// Check row widths
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Create table style
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	cellStyle := lipgloss.NewStyle()

	// Build table
	var table string

	// Headers
	headerRow := ""
	for i, header := range headers {
		paddedHeader := fmt.Sprintf("%-*s", colWidths[i], header)
		headerRow += headerStyle.Render(paddedHeader) + " | "
	}
	headerRow = headerRow[:len(headerRow)-3] // Remove last " | "

	// Separator
	separator := ""
	for i := range headers {
		separator += strings.Repeat("-", colWidths[i]) + " | "
	}
	separator = separator[:len(separator)-3] // Remove last " | "

	// Rows
	var rowStrings []string
	for _, row := range rows {
		rowStr := ""
		for i, cell := range row {
			if i < len(colWidths) {
				paddedCell := fmt.Sprintf("%-*s", colWidths[i], cell)
				rowStr += cellStyle.Render(paddedCell) + " | "
			}
		}
		rowStr = rowStr[:len(rowStr)-3] // Remove last " | "
		rowStrings = append(rowStrings, rowStr)
	}

	// Combine all parts
	table = headerRow + "\n" + separator + "\n"
	for _, row := range rowStrings {
		table += row + "\n"
	}

	fmt.Println(tableStyle.Render(table))
}

// Section creates a section header
func Section(title string) {
	fmt.Println()
	fmt.Println(styles.Title.Render("📋 " + title))
	fmt.Println(styles.Dim.Render(string(make([]byte, len(title)+3, len(title)+3))))
}

// Command shows a command example
func Command(cmd string) {
	fmt.Println(styles.Code.Render("$ " + cmd))
}

// Status shows a status indicator
func Status(status, message string) {
	var icon, style lipgloss.Style

	switch status {
	case "success":
		icon = styles.Success
		style = styles.Success
	case "error":
		icon = styles.Error
		style = styles.Error
	case "warning":
		icon = styles.Warning
		style = styles.Warning
	case "info":
		icon = styles.Info
		style = styles.Info
	default:
		icon = styles.Dim
		style = styles.Dim
	}

	fmt.Printf("%s %s\n", icon.Render("●"), style.Render(message))
}
