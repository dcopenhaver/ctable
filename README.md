# ctable

A Go package for creating beautiful, columnar console table output with automatic formatting and advanced features.

## Features

- **Automatic Column Sizing** - Columns automatically adjust width based on actual data content
- **Smart Truncation** - Control maximum column width with configurable truncation and visual indicators (`...`)
- **Text Justification** - Support for left and right text alignment per column
- **Multiline Cell Values** - Display multiple values within a single cell using string slices
- **UTF-8 Support** - Proper handling of Unicode characters for accurate width calculations
- **Header Display** - Optional column headers with separator lines
- **Zero Dependencies** - Uses only Go standard library packages
- **Type Safety** - Compile-time type checking with structured Column and Table types

## Installation

```bash
go get github.com/dcopenhaver/ctable
```

## Usage

### Basic Example

```go
package main

import "github.com/dcopenhaver/ctable"

func main() {
    // Define columns
    columns := []ctable.Column{
        ctable.NewColumn("Name", 0),      // 0 = no truncation
        ctable.NewColumn("Email", 30),    // truncate at 30 characters
        ctable.NewColumn("Status", 0),
    }

    // Create table
    table := ctable.NewTable(columns)

    // Add rows
    table.AddRow("John Doe", "john.doe@example.com", "Active")
    table.AddRow("Jane Smith", "jane.smith@verylongdomainname.com", "Inactive")
    table.AddRow("Bob Johnson", "bob@example.com", "Active")

    // Display with headers
    table.Display(true)
}
```

**Output:**
```
Name         Email                             Status
===========  ================================  ========
John Doe     john.doe@example.com              Active
Jane Smith   jane.smith@verylongdomainnam...   Inactive
Bob Johnson  bob@example.com                   Active
```

### Multiline Cell Values

```go
columns := []ctable.Column{
    ctable.NewColumn("Product", 0),
    ctable.NewColumn("Features", 0),
    ctable.NewColumn("Price", 0),
}

table := ctable.NewTable(columns)

// Use string slices for multiline values
table.AddRow(
    "Widget Pro",
    []string{"Fast", "Reliable", "Secure", "Scalable"},
    "$99",
)

table.AddRow(
    "Basic Widget",
    "Simple",  // Mix single strings with multiline
    "$49",
)

table.Display(true)
```

**Output:**
```
Product       Features    Price
============  ==========  =====
Widget Pro    Fast        $99
              Reliable
              Secure
              Scalable
Basic Widget  Simple      $49
```

### Custom Justification

```go
columns := []ctable.Column{
    ctable.NewColumn("Item", 0),
    ctable.NewColumn("Quantity", 0),
    ctable.NewColumn("Price", 0),
}

// Set right justification for numeric columns
columns[1].Justification = "right"
columns[2].Justification = "right"

table := ctable.NewTable(columns)
table.AddRow("Apples", "150", "$45.00")
table.AddRow("Bananas", "75", "$22.50")
table.AddRow("Oranges", "200", "$60.00")

table.Display(true)
```

**Output:**
```
Item     Quantity   Price
=======  ========  =======
Apples        150  $45.00
Bananas        75  $22.50
Oranges       200  $60.00
```

### Truncation Control

```go
columns := []ctable.Column{
    ctable.NewColumn("ID", 0),
    ctable.NewColumn("Description", 20),  // Truncate at 20 chars
    ctable.NewColumn("Status", 0),
}

table := ctable.NewTable(columns)
table.AddRow("001", "This is a very long description that will be truncated", "OK")
table.AddRow("002", "Short desc", "OK")

table.Display(true)
```

**Output:**
```
ID   Description             Status
===  ======================  ======
001  This is a very long ...  OK
002  Short desc               OK
```

## API Reference

### Types

#### `Column`
```go
type Column struct {
    Name          string  // Column header name
    Justification string  // "left" or "right" (default: "left")
    // Internal fields managed automatically
}
```

#### `Table`
```go
type Table struct {
    Columns     []Column
    ColumnCount int
    Rows        [][]interface{}
    RowCount    int
}
```

### Functions

#### `NewColumn(name string, truncateAt int) Column`
Creates a new column definition.
- `name`: Column header text
- `truncateAt`: Maximum character width (0 = no truncation)

#### `NewTable(columns []Column) Table`
Creates a new table with the specified column definitions.

#### `(ct *Table) AddRow(fields ...interface{})`
Adds a row to the table. Accepts variadic arguments matching column count.
- Each field can be a `string` or `[]string` (for multiline values)
- Number of fields must match number of columns

#### `(ct *Table) Display(showHeaders bool)`
Renders and prints the table to stdout.
- `showHeaders`: If true, displays column headers with separator line

## Requirements

- Go 1.13 or higher

## License

MIT License

Copyright (c) 2024 David J. Copenhaver

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Disclaimer

This software is provided "as is" without warranty of any kind, either express or implied. The author and contributors are not responsible for any damages, data loss, or other issues that may arise from using this software. Use at your own risk.

This package is intended for console/terminal output formatting only. It is not designed for production-critical systems or applications where output formatting errors could cause significant harm.

Users are responsible for:
- Testing the package thoroughly in their specific use cases
- Validating output formatting meets their requirements
- Ensuring compatibility with their target environments and terminal emulators
- Handling any errors or unexpected behavior appropriately in their applications

## Author

David J. Copenhaver

**Note:** This is a personal tool shared publicly for accessibility. Issues and pull requests are disabled as this is not actively maintained as an open-source project. Feel free to fork and modify for your own use.
