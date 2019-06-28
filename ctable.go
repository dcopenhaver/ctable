package ctable

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ConsoleTableColumn struct {
	Name                string
	truncate_at         int
	Justification       string
	truncation_required bool
	max_length          int
}

func NewConsoleTableColumn(name string, truncate_at int) ConsoleTableColumn {

	/*
		truncation and max length related items need to inlude the column name value as well,
		because the column names are a part of the data set when it comes to display logic
	*/

	return ConsoleTableColumn{
		Name:                name,
		truncate_at:         truncate_at,
		Justification:       "left",
		truncation_required: truncate_at > 0 && utf8.RuneCountInString(name) > truncate_at,
		max_length:          utf8.RuneCountInString(name),
	}
}

type ConsoleTable struct {
	Columns     []ConsoleTableColumn
	ColumnCount int
	Rows        [][]string
	RowCount    int
}

func NewConsoleTable(columns []ConsoleTableColumn) ConsoleTable {
	return ConsoleTable{
		Columns:     columns,
		ColumnCount: len(columns),
		Rows:        [][]string{},
		RowCount:    0,
	}
}

func (ct *ConsoleTable) AddRow(fields ...string) {

	if len(fields) != ct.ColumnCount {
		log.Fatal("CONSOLETABLE: Cannot add a row of data with more, or fewer, fields than defined columns.")
	}

	// update max length values and truncation status stored with column defs.
	// whether truncation *will* be required is stored with the column def so it can be used in display logic,
	// but we only truncate upon display so we keep all the data
	// note: max length is initialized to the column name's length when new column object is instantiated,
	// as it is essentially a part of the data set when it comes to display logic

	for i := 0; i < ct.ColumnCount; i++ {
		if ct.Columns[i].max_length < utf8.RuneCountInString(fields[i]) {
			ct.Columns[i].max_length = utf8.RuneCountInString(fields[i])
		}
		if ct.Columns[i].truncate_at > 0 && ct.Columns[i].max_length > ct.Columns[i].truncate_at {
			ct.Columns[i].truncation_required = true
		}
	}

	ct.Rows = append(ct.Rows, fields)
}

func (ct *ConsoleTable) Display(show_headers bool) {

	processed_rows := []string{}

	for _, row := range ct.Rows {
		// for each row
		row_str := ""

		for i, col := range ct.Columns {
			// for each field - build row string including padding for columnar output, justification, and any truncation per column defs
			field_data := row[i]

			// truncate field value?
			if col.truncation_required && utf8.RuneCountInString(field_data) > col.truncate_at {
				field_data = field_data[:col.truncate_at] + "..."
			}

			// create format string that will be used for column width and justification
			var format_string string
			var just_code string // used inside format string

			if col.Justification == "left" {
				just_code = "%-"
			} else {
				just_code = "%"
			}

			if col.truncation_required {
				format_string = just_code + strconv.Itoa(col.truncate_at+3) + "v" // +3 for the ... added when truncated
			} else {
				format_string = just_code + strconv.Itoa(col.max_length) + "v"
			}

			// padding between columns, prepend a space to all but the first column
			if i == 0 {
				row_str += fmt.Sprintf(format_string, field_data)
			} else {
				row_str += " " + fmt.Sprintf(format_string, field_data)
			}
		} // END for each column

		processed_rows = append(processed_rows, row_str)
	} // END for each row

	if show_headers {
		header_str := ""
		header_separator := ""

		for i := range ct.Columns { // note: this format of 'range' is pointer rather than value, as we need to modify original object

			col := ct.Columns[i] // <- still pointer format for range, but assigning to var as it was already used throughout

			// did we truncate? if so header needs to account for that
			if col.truncation_required {

				// padding between columns, prepend space to all but first column
				if i == 0 {
					header_separator += strings.Repeat("=", col.truncate_at+3) // +3 to account for the '...'
				} else {
					header_separator += " " + strings.Repeat("=", col.truncate_at+3) // +3 to account for the '...'
				}

				// truncate column name also?
				if utf8.RuneCountInString(col.Name) > col.truncate_at {
					// padding between columns, prepend space to all but first column
					if i == 0 {
						header_str += fmt.Sprintf("%-"+strconv.Itoa(col.truncate_at+3)+"v", col.Name[:col.truncate_at]+"...")
					} else {
						header_str += " " + fmt.Sprintf("%-"+strconv.Itoa(col.truncate_at+3)+"v", col.Name[:col.truncate_at]+"...")
					}
				} else {
					// padding between columns, prepend a space to all but first column
					if i == 0 {
						header_str += fmt.Sprintf("%-"+strconv.Itoa(col.truncate_at+3)+"v", col.Name)
					} else {
						header_str += " " + fmt.Sprintf("%-"+strconv.Itoa(col.truncate_at+3)+"v", col.Name)
					}
				}

			} else {
				// padding between columns, prepend a space to all but the first column
				if i == 0 {
					header_str += fmt.Sprintf("%-"+strconv.Itoa(col.max_length)+"v", col.Name)
					header_separator += strings.Repeat("=", col.max_length)
				} else {
					header_str += " " + fmt.Sprintf("%-"+strconv.Itoa(col.max_length)+"v", col.Name)
					header_separator += " " + strings.Repeat("=", col.max_length)
				}
			}
		}
		// output header
		fmt.Println(header_str)
		fmt.Println(header_separator)
	}

	for _, r := range processed_rows {
		fmt.Println(r)
	}
}
