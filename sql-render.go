package main

import (
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"strings"
)

func highlightSQL(input string) string {
	keywords := map[string]lipgloss.Style{
		"SELECT":   lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"FROM":     lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"WHERE":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"GROUP BY": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"ORDER BY": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"HAVING":   lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"LIMIT":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"OFFSET":   lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"DISTINCT": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),

		"INSERT": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"INTO":   lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"VALUES": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"UPDATE": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"SET":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"DELETE": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),

		"CREATE":   lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"TABLE":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"VIEW":     lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"INDEX":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"ALTER":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"DROP":     lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"TRUNCATE": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),

		"INT":       lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"VARCHAR":   lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"TEXT":      lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"BOOLEAN":   lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"DATE":      lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"TIMESTAMP": lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),
		"SERIAL":    lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(false),

		"JOIN":  lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"INNER": lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"OUTER": lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"LEFT":  lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"RIGHT": lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"FULL":  lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"ON":    lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),

		"AND":     lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"OR":      lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"NOT":     lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"IN":      lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"LIKE":    lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"BETWEEN": lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"IS":      lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"NULL":    lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),

		"COUNT": lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(false),
		"SUM":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(false),
		"AVG":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(false),
		"MIN":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(false),
		"MAX":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(false),

		"BEGIN":       lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		"COMMIT":      lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		"ROLLBACK":    lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		"TRANSACTION": lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),

		"AS":         lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"EXISTS":     lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true),
		"UNION":      lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"ALL":        lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"DEFAULT":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"PRIMARY":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"KEY":        lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"FOREIGN":    lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"REFERENCES": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"CONSTRAINT": lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"CHECK":      lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"UNIQUE":     lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(false),
		"WITH":       lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
	}

	var result strings.Builder
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	})

	space := regexp.MustCompile(`\s+`)
	spaces := space.FindAllString(input, -1)

	for i, word := range words {
		upperWord := strings.ToUpper(word)
		if style, exists := keywords[upperWord]; exists {
			result.WriteString(style.Render(word))
		} else {
			result.WriteString(word)
		}

		if i < len(spaces) {
			result.WriteString(spaces[i])
		}
	}

	return result.String()
}

func highlightForEditor(input string) string {
	cursorPrefix := ""
	cursorSuffix := ""

	if strings.HasPrefix(input, "\x1b[?25") {
		cursorPrefix = input[:4]
		input = input[4:]
	}

	if strings.HasSuffix(input, "\x1b[?25") {
		cursorSuffix = input[len(input)-4:]
		input = input[:len(input)-4]
	}

	highlighted := highlightSQL(input)

	return cursorPrefix + highlighted + cursorSuffix
}
