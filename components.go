package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func setupTextarea() textarea.Model {
	editor := textarea.New()
	editor.Placeholder = "SELECT * FROM users;"
	editor.Focus()
	editor.Prompt = "┃ "
	editor.CharLimit = 10000
	editor.FocusedStyle.Base = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("235"))
	editor.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("236"))
	editor.Cursor.Blink = true
	editor.Cursor.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("15")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	return editor
}

func setupTable() table.Model {
	columns := []table.Column{}
	rows := []table.Row{}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(1),
		table.WithStyles(table.Styles{
			Cell:     lipgloss.NewStyle().Padding(0, 1),
			Header:   lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true),
			Selected: lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false),
		}),
	)

	if len(columns) > 0 {
		t.SetWidth(80)
	}

	return t
}
