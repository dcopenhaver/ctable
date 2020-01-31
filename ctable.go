package ctable

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Column struct {
	Name               string
	truncateAt         int
	Justification      string
	truncationRequired bool
	maxLength          int
}

func NewColumn(name string, truncateAt int) Column {

	/*
		truncation and max length related items need to inlude the column name value as well,
		because the column names are a part of the data set when it comes to display logic
	*/

	return Column{
		Name:               name,
		truncateAt:         truncateAt,
		Justification:      "left",
		truncationRequired: truncateAt > 0 && utf8.RuneCountInString(name) > truncateAt,
		maxLength:          utf8.RuneCountInString(name),
	}
}

type Table struct {
	Columns     []Column
	ColumnCount int
	Rows        [][]string
	RowCount    int
}

func NewTable(columns []Column) Table {
	return Table{
		Columns:     columns,
		ColumnCount: len(columns),
		Rows:        [][]string{},
		RowCount:    0,
	}
}

func (t *Table) AddRow(fields ...string) {

	if len(fields) != t.ColumnCount {
		log.Fatal("CONSOLETABLE: Cannot add a row of data with more, or fewer, fields than defined columns.")
	}

	// update max length values and truncation status stored with column defs.
	// whether truncation *will* be required is stored with the column def so it can be used in display logic,
	// but we only truncate upon display so we keep all the data
	// note: max length is initialized to the column name's length when new column object is instantiated,
	// as it is essentially a part of the data set when it comes to display logic

	for i := 0; i < t.ColumnCount; i++ {
		if t.Columns[i].maxLength < utf8.RuneCountInString(fields[i]) {
			t.Columns[i].maxLength = utf8.RuneCountInString(fields[i])
		}
		if t.Columns[i].truncateAt > 0 && t.Columns[i].maxLength > t.Columns[i].truncateAt {
			t.Columns[i].truncationRequired = true
		}
	}

	t.Rows = append(t.Rows, fields)
}

func (t *Table) Display(showHeaders bool) {

	processedRows := []string{}

	for _, row := range t.Rows {
		// for each row
		rowStr := ""

		for i, col := range t.Columns {
			// for each field - build row string including padding for columnar output, justification, and any truncation per column defs
			fieldData := row[i]

			// truncate field value?
			if col.truncationRequired && utf8.RuneCountInString(fieldData) > col.truncateAt {
				fieldData = fieldData[:col.truncateAt] + "..."
			}

			// create format string that will be used for column width and justification
			var formatString string
			var justCode string // used inside format string

			if col.Justification == "left" {
				justCode = "%-"
			} else {
				justCode = "%"
			}

			if col.truncationRequired {
				formatString = justCode + strconv.Itoa(col.truncateAt+3) + "v" // +3 for the ... added when truncated
			} else {
				formatString = justCode + strconv.Itoa(col.maxLength) + "v"
			}

			// padding between columns, prepend a space to all but the first column
			if i == 0 {
				rowStr += fmt.Sprintf(formatString, fieldData)
			} else {
				rowStr += " " + fmt.Sprintf(formatString, fieldData)
			}
		} // END for each column

		processedRows = append(processedRows, rowStr)
	} // END for each row

	if showHeaders {
		headerStr := ""
		headerSeparator := ""

		for i := range t.Columns { // note: this format of 'range' is pointer rather than value, as we need to modify original object

			col := t.Columns[i] // <- still pointer format for range, but assigning to var as it was already used throughout

			// did we truncate? if so header needs to account for that
			if col.truncationRequired {

				// padding between columns, prepend space to all but first column
				if i == 0 {
					headerSeparator += strings.Repeat("=", col.truncateAt+3) // +3 to account for the '...'
				} else {
					headerSeparator += " " + strings.Repeat("=", col.truncateAt+3) // +3 to account for the '...'
				}

				// truncate column name also?
				if utf8.RuneCountInString(col.Name) > col.truncateAt {
					// padding between columns, prepend space to all but first column
					if i == 0 {
						headerStr += fmt.Sprintf("%-"+strconv.Itoa(col.truncateAt+3)+"v", col.Name[:col.truncateAt]+"...")
					} else {
						headerStr += " " + fmt.Sprintf("%-"+strconv.Itoa(col.truncateAt+3)+"v", col.Name[:col.truncateAt]+"...")
					}
				} else {
					// padding between columns, prepend a space to all but first column
					if i == 0 {
						headerStr += fmt.Sprintf("%-"+strconv.Itoa(col.truncateAt+3)+"v", col.Name)
					} else {
						headerStr += " " + fmt.Sprintf("%-"+strconv.Itoa(col.truncateAt+3)+"v", col.Name)
					}
				}

			} else {
				// padding between columns, prepend a space to all but the first column
				if i == 0 {
					headerStr += fmt.Sprintf("%-"+strconv.Itoa(col.maxLength)+"v", col.Name)
					headerSeparator += strings.Repeat("=", col.maxLength)
				} else {
					headerStr += " " + fmt.Sprintf("%-"+strconv.Itoa(col.maxLength)+"v", col.Name)
					headerSeparator += " " + strings.Repeat("=", col.maxLength)
				}
			}
		}
		// output header
		fmt.Println(headerStr)
		fmt.Println(headerSeparator)
	}

	for _, r := range processedRows {
		fmt.Println(r)
	}
}
