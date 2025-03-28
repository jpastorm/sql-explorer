package main

import (
	"database/sql"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/muesli/termenv"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
)

const (
	focusEditor = iota
	focusList
	focusResults
)

type dbItem struct {
	name  string
	kind  string // "db", "table", "column"
	child []dbItem
}

func (i dbItem) Title() string       { return i.name }
func (i dbItem) Description() string { return i.kind }
func (i dbItem) FilterValue() string { return i.name }

type model struct {
	dbList        list.Model
	editor        textarea.Model
	db            *sql.DB
	data          []dbItem // Loaded databases
	insideColumns bool
	resultWindow  bool // Indicates if the results window is displayed
	queryError    string
	queryResult   []string // Query results
	currentPage   int      // Current page of the results
	itemsPerPage  int      // Number of rows per page
	focusedEditor bool     // Indicates if the focus is on the editor
	currentTable  string   // Name of the current table

	resultsTable table.Model
	showResults  bool
	focusState   int

	LWidth     int
	EWidth     int
	RHeight    int
	MainHeight int
	TotalWidth int
}

func initialModel() model {
	db, err := connectToPostgres()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]list.Item, len(dbs))
	for i, db := range dbs {
		items[i] = db
	}

	InitBackupSystems()
	editor := setupTextarea()
	backup, err := loadEditorBackup()
	if err == nil {
		editor.SetValue(backup)
	}

	tbl := setupTable()
	return model{
		dbList:       list.New(items, list.NewDefaultDelegate(), 0, 0),
		editor:       editor,
		db:           db,
		data:         dbs,
		itemsPerPage: 10,
		resultsTable: tbl,
		focusState:   focusEditor,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TotalWidth = int(float64(msg.Width) * 0.95)
		m.LWidth = int(float64(m.TotalWidth) * 0.30)
		m.EWidth = m.TotalWidth - m.LWidth - 4

		totalHeight := msg.Height - 6
		m.MainHeight = int(float64(totalHeight) * 0.60)
		m.RHeight = totalHeight - m.MainHeight - 3

		m.reDrawTable()

		m.dbList.SetSize(m.LWidth, m.MainHeight-4)
		m.editor.SetWidth(m.EWidth)
		m.editor.SetHeight(m.MainHeight)
	}

	if m.showResults && m.focusState == focusResults {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.showResults = false
				m.focusState = focusEditor
			case "tab":
				m.focusState = focusEditor
				m.resultsTable.Blur()
				m.editor.Focus()
				m.dbList.SetFilteringEnabled(false)
			default:
				m.resultsTable, cmd = m.resultsTable.Update(msg)
				return m, cmd
			}
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.queryError = ""
			m.resultWindow = false
		case "tab":
			switch m.focusState {
			case focusEditor:
				m.focusState = focusList
				m.editor.Blur()
				m.dbList.SetFilteringEnabled(true)
			case focusList:
				if m.showResults {
					m.focusState = focusResults
					m.dbList.SetFilteringEnabled(false)
					m.resultsTable.Focus()
				} else {
					m.focusState = focusEditor
					m.dbList.SetFilteringEnabled(false)
					m.editor.Focus()
				}
			case focusResults:
				m.focusState = focusEditor
				m.resultsTable.Blur()
				m.editor.Focus()
			}
		case "backspace":
			if m.focusState != focusEditor && m.insideColumns {
				m.currentTable = ""
				tables, _ := getTables(m.db)
				var items []list.Item
				for _, table := range tables {
					items = append(items, list.Item(table))
				}
				m.dbList.SetItems(items)
				m.insideColumns = false
			}
		case "enter":
			if m.focusState != focusEditor {
				selectedItem := m.dbList.SelectedItem().(dbItem)
				if selectedItem.kind == "tables" {
					m.currentTable = selectedItem.name
					columns, _ := getColumns(m.db, selectedItem.name)
					selectedItem.child = columns
				}

				if len(selectedItem.child) > 0 {
					var items []list.Item
					for _, child := range selectedItem.child {
						items = append(items, list.Item(child))
					}
					m.dbList.SetItems(items)
					m.insideColumns = true
				}
			}
		//case "ctrl+v":
		//	if m.focusedEditor {
		//		text, err := clipboard.ReadAll()
		//		if err != nil {
		//			log.Fatalf("Error reading clipboard: %v", err)
		//		}
		//		m.editor.InsertString(text)
		//	}
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+c":
			if m.focusState == focusEditor {
				currentQuery := extractCurrentLine(m.editor)
				currentQuery = strings.TrimSpace(currentQuery)
				err := clipboard.WriteAll(currentQuery)
				if err != nil {
					log.Fatalf("Error clipboard process: %v", err)
				}
			}
		case "ctrl+a":
			if m.focusState == focusEditor {
				err := clipboard.WriteAll(m.editor.Value())
				if err != nil {
					log.Fatalf("Error clipboard process: %v", err)
				}
			} else {
				return m, tea.Quit
			}
		case "ctrl+x":
			if m.focusState == focusEditor {
				content := m.editor.Value()
				if content == "" {
					return m, nil
				}
				lines := strings.Split(content, "\n")
				currentLine := m.editor.Line()
				if currentLine >= len(lines) {
					currentLine = len(lines) - 1
				}

				clipboard.WriteAll(lines[currentLine])
				if len(lines) == 1 {
					m.editor.SetValue("")
				} else {
					lines = append(lines[:currentLine], lines[currentLine+1:]...)
					m.editor.SetValue(strings.Join(lines, "\n"))
					newPos := currentLine
					if newPos >= len(lines) {
						newPos = len(lines) - 1
					}
					if newPos > 0 && lines[newPos] == "" {
						newPos--
					}
				}
			}
		case "ctrl+y":
			if m.focusState == focusEditor {
				currentQuery := extractCurrentLine(m.editor)
				currentQuery = strings.TrimSpace(currentQuery)
				if currentQuery != "" {
					rows, err := m.db.Query(currentQuery)
					if err != nil {
						m.queryError = err.Error()
						break
					}
					defer rows.Close()

					columns, _ := rows.Columns()
					tableColumns := make([]table.Column, len(columns))
					for i, col := range columns {
						tableColumns[i] = table.Column{
							Title: col,
							Width: len(col) + 2,
						}
					}

					m.resultsTable = table.New(
						table.WithColumns(tableColumns),
						table.WithRows([]table.Row{}),
						table.WithWidth(m.EWidth),
						table.WithHeight(m.RHeight),
					)

					var tableRows []table.Row
					for rows.Next() {
						values := make([]interface{}, len(columns))
						valuePtrs := make([]interface{}, len(columns))
						for i := range values {
							valuePtrs[i] = &values[i]
						}
						if err := rows.Scan(valuePtrs...); err != nil {
							log.Println("Error scan:", err)
							break
						}

						row := make([]string, len(values))
						for i, val := range values {
							row[i] = fmt.Sprintf("%v", val)
						}
						tableRows = append(tableRows, row)
					}

					m.resultsTable.SetColumns(tableColumns)
					m.resultsTable.SetRows(tableRows)
					m.showResults = true
					m.queryError = ""
					m.currentPage = 0
					m.updateResultsTable(tableRows)
					SaveTableState(m.resultsTable)
				}
			}
		}
	}

	switch m.focusState {
	case focusEditor:
		m.editor, cmd = m.editor.Update(msg)
	case focusList:
		m.dbList, cmd = m.dbList.Update(msg)
	case focusResults:
		m.resultsTable, cmd = m.resultsTable.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	m.reDrawTable()
	totalWidth := m.LWidth + m.EWidth + 4
	containerStyle := lipgloss.NewStyle().
		Width(totalWidth).
		MarginLeft(2).
		MarginRight(2)

	listStyle := lipgloss.NewStyle().
		Width(m.LWidth).
		MaxHeight(m.MainHeight+10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	editorStyle := lipgloss.NewStyle().
		Width(m.EWidth).
		Height(m.MainHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	switch m.focusState {
	case focusEditor:
		editorStyle = editorStyle.
			BorderForeground(lipgloss.Color("5")).
			Background(lipgloss.Color("235"))
	case focusList:
		listStyle = listStyle.
			BorderForeground(lipgloss.Color("5")).
			Background(lipgloss.Color("235"))
	}

	mainSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listStyle.Render(m.dbList.View()),
		editorStyle.Render(highlightForEditor(m.editor.View())),
	)

	resultsStyle := lipgloss.NewStyle().
		Width(m.TotalWidth-4).
		Height(m.RHeight-2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 0).
		MarginTop(1)

	tableContentStyle := lipgloss.NewStyle().
		Width(m.TotalWidth - 6).
		Height(m.RHeight - 4)

	resultsContent := ""
	if m.showResults {
		resultsContent = tableContentStyle.Render(m.resultsTable.View())
	}

	if m.focusState == focusResults {
		resultsStyle = resultsStyle.
			BorderForeground(lipgloss.Color("5")).
			Background(lipgloss.Color("235"))
	}

	resultsSection := resultsStyle.Render(resultsContent)

	statusBar := ""
	if m.queryError != "" {
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("196")).
			Bold(true).
			Padding(0, 1).
			Width(totalWidth).
			Render("Error: " + m.queryError)
	} else if m.currentTable != "" {
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Width(totalWidth).
			Render("Current Table: " + m.currentTable)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		containerStyle.Render(mainSection),
		lipgloss.NewStyle().MarginTop(1).Render(resultsSection),
		lipgloss.NewStyle().MarginTop(1).Render(statusBar),
	)
}

func main() {
	defer CloseBackupSystems()
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	clearScreen()
	lipgloss.SetColorProfile(termenv.TrueColor)
	p := tea.NewProgram(initialModel())

	finalModelChan := make(chan tea.Model, 1)
	go func() {
		m, err := p.Run()
		if err != nil {
			log.Println("Runtime error:", err)
		}
		finalModelChan <- m
	}()

	finalModel := <-finalModelChan

	if m, ok := finalModel.(model); ok {
		err := saveEditorBackup(m.editor.Value())
		if err != nil {
			log.Println("Final backup save error:", err)
		}
	}
}
