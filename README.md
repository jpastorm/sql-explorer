# SQL Explorer

![SQL Explorer Demo](demo.gif) *Replace with actual demo GIF*

## Overview

SQL Explorer designed terminal user interface (TUI) for interacting with PostgreSQL databases. With its modern aesthetic and intuitive controls, it provides a seamless database exploration experience right in your terminal.

## Features

- **Database Navigation**: Intuitive tree-like navigation through databases, tables, and columns
- **Query Execution**: Run SQL queries with immediate results
- **Clipboard Integration**: Copy queries and results with simple keyboard shortcuts
- **Result Pagination**: View large result sets with built-in pagination
- **Query History**: Automatic backup of your query history
- **Responsive Layout**: Adapts to different terminal sizes
- **Keyboard-Centric**: Designed for efficient keyboard navigation

## Installation

### Prerequisites
- Go 1.16+
- PostgreSQL client libraries

### Build from Source
```bash
git clone https://github.com/yourusername/sql-explorer.git
cd sql-explorer
go build -o sqlexplorer
```
# Configuration

## Using .env file

Create a `.env` file in the project root with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=postgres
DB_SSLMODE=disable
```

The application will automatically load these connection parameters. You can also override them with environment variables:

```bash
export DB_HOST=myserver
export DB_PORT=5433
./sqlexplorer
```

## Usage

```bash
./sqlexplorer
```

## Key Bindings

```
Key Combination   Action
----------------  ----------------------------------------------
Tab              Cycle focus between editor, database list, and results
Ctrl+y           Execute current query
Ctrl+c           Copy current line to clipboard
Ctrl+a           Copy entire query to clipboard
Ctrl+x           Cut current line
Ctrl+q           Quit application
Enter            Navigate into table/column
Backspace        Navigate back from columns view
Esc              Clear error messages or exit results view
```

## Technical Details

### Built With

- **BubbleTea** - TUI framework
- **LipGloss** - Style definitions
- **PostgreSQL driver** - Database connectivity
- **godotenv** - Environment variables loader

### Architecture

The application follows the **Model-View-Update** pattern implemented by BubbleTea, with separate components for:

- Database navigation tree
- SQL query editor
- Results display
- Configuration loader

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License
