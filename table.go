package main

import (
	"github.com/charmbracelet/bubbles/table"
	"unicode"
)

const (
	maxRowsToRender = 1000
)

func (m *model) calculateColumnWidths(cols []table.Column, rows []table.Row, availableWidth, minColWidth, maxColWidth int) []int {
	if len(rows) == 0 || len(cols) == 0 {
		colWidth := availableWidth / len(cols)
		result := make([]int, len(cols))
		for i := range result {
			result[i] = clamp(colWidth, min(len(cols[i].Title), availableWidth), maxColWidth)
		}
		return result
	}

	for i := range cols {
		minColWidth = max(minColWidth, len(cols[i].Title))
	}

	priorities := m.prioritizeColumns(cols, rows)
	totalPriority := sum(priorities...)

	colWidths := make([]int, len(cols))
	remainingWidth := availableWidth

	for i := range cols {
		if totalPriority > 0 {
			colWidths[i] = (priorities[i] * availableWidth) / totalPriority
		} else {
			colWidths[i] = availableWidth / len(cols)
		}
		colWidths[i] = clamp(colWidths[i], len(cols[i].Title), maxColWidth)
		remainingWidth -= colWidths[i]
	}

	for remainingWidth > 0 {
		distributed := false

		for i := range cols {
			if colWidths[i] < maxColWidth && float64(colWidths[i]) < float64(priorities[i])*1.2 {
				colWidths[i]++
				remainingWidth--
				distributed = true
				if remainingWidth == 0 {
					break
				}
			}
		}

		if remainingWidth > 0 && !distributed {
			for i := range cols {
				if colWidths[i] < maxColWidth {
					colWidths[i]++
					remainingWidth--
					if remainingWidth == 0 {
						break
					}
				}
			}
		}
	}

	return colWidths
}

func (m *model) prioritizeColumns(cols []table.Column, rows []table.Row) []int {
	priorities := make([]int, len(cols))

	for i := range cols {
		maxLen := max(len(cols[i].Title), 5)
		uniqueValues := make(map[string]bool)
		totalContentLen := 0
		isNumericID := true
		maxContentLen := 0

		for _, row := range rows {
			cellValue := row[i]
			cellLen := len(cellValue)

			if cellLen > maxLen {
				maxLen = cellLen
			}
			if cellLen > maxContentLen {
				maxContentLen = cellLen
			}

			if isNumericID {
				for _, r := range cellValue {
					if !unicode.IsDigit(r) {
						isNumericID = false
						break
					}
				}
			}

			uniqueValues[cellValue] = true
			totalContentLen += cellLen
		}

		uniquenessFactor := len(uniqueValues)
		avgContentLen := 0
		if len(rows) > 0 {
			avgContentLen = totalContentLen / len(rows)
		}

		switch {
		case isNumericID:
			priorities[i] = max(len(cols[i].Title), 5)
		case maxContentLen > 20:
			priorities[i] = maxLen * (uniquenessFactor + 3) * (avgContentLen + 2)
		default:
			priorities[i] = maxLen * (uniquenessFactor + 2) * (avgContentLen + 1)
		}
	}

	return priorities
}

func (m model) updateResultsTable(rows []table.Row) {
	if len(rows) > maxRowsToRender {
		rows = rows[:maxRowsToRender]
	}

	m.resultsTable.SetRows(rows)

	if len(rows) > 0 {
		cols := m.resultsTable.Columns()
		availableWidth := m.TotalWidth - 6

		minColWidth := 4
		maxColWidth := 50

		colWidths := m.calculateColumnWidths(cols, rows, availableWidth, minColWidth, maxColWidth)

		for i := range cols {
			cols[i].Width = colWidths[i]
		}
		m.resultsTable.SetColumns(cols)
	}
}

func (m *model) reDrawTable() {
	tableWidth := m.TotalWidth - 4
	tableHeight := m.RHeight - 2
	m.resultsTable.SetWidth(tableWidth)
	m.resultsTable.SetHeight(tableHeight)
}

func sum(values ...int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
