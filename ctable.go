package ctable

/*
Purpose: To create nice columnar output (tables) for console apps.
Author: David J. Copenhaver

Notable features:
- Automatic column sizing based on the actual data
- Control over truncation
- Control over justification
- Support for mulitline values
*/

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
	Rows        [][]interface{}
	RowCount    int
}

func NewTable(columns []Column) Table {
	return Table{
		Columns:     columns,
		ColumnCount: len(columns),
		Rows:        [][]interface{}{},
		RowCount:    0,
	}
}

/*
	New version of AddRow() to support multiline values in fields - original/previous version commented out below this one

	Example calls:
	ct.AddRow("string data", "string data", "string data")
	ct.AddRow("string data", []string{"one", "two", "three", "four"}, "string data")

	Its variadic, so any number of values in any mix of strings and string slices can be used (number of args has to match number of columns of course)
*/
func (ct *Table) AddRow(fields ...interface{}) {

	if len(fields) != ct.ColumnCount {
		log.Fatal("CONSOLETABLE: Cannot add a row of data with more, or fewer, fields than defined columns.")
	}

	/*
		Update max length values and truncation status stored with column defs.
		Whether truncation *will* be required is stored with the column def so it can be used in display logic,
		but we only truncate upon display so we keep all the data.

		Note: max length is initialized to the column name's length when new column object is instantiated,
		as it is essentially a part of the data set when it comes to display logic.
	*/

	rowState := struct {
		hasMultilineValue bool
		mlFieldIndexes    []int       // multiline field indexes (the field indexes that contain a list of values instead of just one, a slice)
		mlTracker         map[int]int // multiline tracker (need to track which list item to display next for each multiline field)
	}{}

	rowState.hasMultilineValue = false
	rowState.mlTracker = make(map[int]int)

	for i := 0; i < ct.ColumnCount; i++ {

		switch v := fields[i].(type) {

		case string:

			if ct.Columns[i].maxLength < utf8.RuneCountInString(v) {
				ct.Columns[i].maxLength = utf8.RuneCountInString(v)
			}
			if ct.Columns[i].truncateAt > 0 && ct.Columns[i].maxLength > ct.Columns[i].truncateAt {
				ct.Columns[i].truncationRequired = true
			}

		case []string:

			rowState.hasMultilineValue = true
			rowState.mlFieldIndexes = append(rowState.mlFieldIndexes, i)
			rowState.mlTracker[i] = 0

			for _, str := range v {

				if ct.Columns[i].maxLength < utf8.RuneCountInString(str) {
					ct.Columns[i].maxLength = utf8.RuneCountInString(str)
				}
				if ct.Columns[i].truncateAt > 0 && ct.Columns[i].maxLength > ct.Columns[i].truncateAt {
					ct.Columns[i].truncationRequired = true
				}
			}

		default:
			log.Fatal("CONSOLETABLE: You can add only string or []string types as individual fields to AddRow().")
		}
	}

	if !rowState.hasMultilineValue {
		// add as normal, each field is an interface{} that's value IS a single string
		ct.Rows = append(ct.Rows, fields)
	} else {
		// deal with multiline fields

		// need to know the highest number of mulitline values among all multiline fields in order to control right number of iterations to handle them all
		longestMulti := 0

		for _, fi := range rowState.mlFieldIndexes {
			if len(fields[fi].([]string)) > longestMulti {
				longestMulti = len(fields[fi].([]string))
			}
		}

		for x := 0; x < longestMulti; x++ { // in context of the one 'row', this is the 'down' direction due to multiline values

			// contruct slice to be used to append to ct.Rows
			var tempFields []interface{}
			for fi := 0; fi < ct.ColumnCount; fi++ { // ... and this is the 'across' direction

				var isMLField bool
				for _, mfi := range rowState.mlFieldIndexes {
					if fi == mfi {
						isMLField = true
					}
				}

				// first row special
				if x == 0 {
					if !isMLField {
						tempFields = append(tempFields, fields[fi])
					} else {
						tempFields = append(tempFields, fields[fi].([]string)[0])
						rowState.mlTracker[fi]++
					}
				} else {
					// not first row, the rest (to support multiline in the table output, empty strings have to be inserted into all the other fields that are not multiline, *after* the first line)
					if !isMLField {
						tempFields = append(tempFields, "") // blanks in cells/fields next to cells/fields with multilines after first value is displayed with the rest of the non multiline value row
					} else {
						tempFields = append(tempFields, fields[fi].([]string)[x])
						rowState.mlTracker[fi]++ // keeping track of which item in the multiline list has been handled and which is up next
					}
				}
			}

			//  tempFields should have been appropriately built now to add as a 'row' in the context of a row to be displayed across the screen
			ct.Rows = append(ct.Rows, tempFields)
		}
	}
}

//func (t *Table) AddRow(fields ...string) {
//
//	if len(fields) != ct.ColumnCount {
//		log.Fatal("CONSOLETABLE: Cannot add a row of data with more, or fewer, fields than defined columns.")
//	}
//
//	// update max length values and truncation status stored with column defs.
//	// whether truncation *will* be required is stored with the column def so it can be used in display logic,
//	// but we only truncate upon display so we keep all the data
//	// note: max length is initialized to the column name's length when new column object is instantiated,
//	// as it is essentially a part of the data set when it comes to display logic
//
//	for i := 0; i < ct.ColumnCount; i++ {
//		if ct.Columns[i].maxLength < utf8.RuneCountInString(fields[i]) {
//			ct.Columns[i].maxLength = utf8.RuneCountInString(fields[i])
//		}
//		if ct.Columns[i].truncateAt > 0 && ct.Columns[i].maxLength > ct.Columns[i].truncateAt {
//			ct.Columns[i].truncationRequired = true
//		}
//	}
//
//	ct.Rows = append(ct.Rows, fields)
//}

func (ct *Table) Display(showHeaders bool) {

	processedRows := []string{}

	for _, row := range ct.Rows {
		// for each row
		rowStr := ""

		for i, col := range ct.Columns {
			// for each field - build row string including padding for columnar output, justification, and any truncation per column defs
			fieldData := row[i].(string) // to support multiline values, interfaces were used, at this point each field item should be a single string

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

		for i := range ct.Columns { // note: this format of 'range' is pointer rather than value, as we need to modify original object

			col := ct.Columns[i] // <- still pointer format for range, but assigning to var as it was already used throughout

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
