# simple-jot
A CLI tool for taking and managing notes (written in Go)

## Installation

```bash
go build -o simple-jot .
```

## Usage

### Important Note About Quotes
When providing note content through command-line flags, use single quotes (`'`) instead of double quotes (`"`) to avoid shell interpretation issues:

```bash
# Good - uses single quotes
simple-jot create "My Note" -n 'This is my note content!'

# Problematic - double quotes might cause issues with special characters
simple-jot create "My Note" -n "This is my note content!"
```

### Commands

#### Create Notes
Create a new note with content:
```bash
# Create with direct content
simple-jot create "Note Title" -n 'Your note content'

# Create and set as active note
simple-jot create "Note Title" -n 'Your note content' -s

# Create by piping content from a file
cat file.txt | simple-jot create "Note Title"

# Create by piping content directly
echo 'Note content' | simple-jot create "Note Title"
```

#### Edit Notes
Edit an existing note:
```bash
# Overwrite content
simple-jot edit <note-id> -n 'New content'

# Append content
simple-jot edit <note-id> -a 'Additional content'

# Edit by piping new content
echo 'New content' | simple-jot edit <note-id>
```

#### List Notes
View all notes:
```bash
simple-jot list
```

#### Search Notes
Search through your notes:
```bash
# Search by date range
simple-jot search --date-start "2025-01-01" --date-end "2025-02-01"

# Search by content
simple-jot search --content "search term"

# Search by tag
simple-jot search --tag "important"

# Semantic search
simple-jot search --semantic "programming"
```

#### Tag Notes
Add tags to your notes:
```bash
# Tag a note
simple-jot tag <note-id> "tag-name"
```

#### Delete Notes
Delete a note:
```bash
simple-jot delete <note-id>
```

#### Configuration
Manage your configuration:
```bash
# Set active note
simple-jot config set note <note-id>

# Get active note
simple-jot config get note
```

## Features
- Create and manage notes with unique IDs
- Edit notes with overwrite or append functionality
- Search notes by content, date, or tags
- Tag system for organization
- Pipe content from files or other commands
- Configuration management
- Table-formatted output for better readability

## Tips
1. Use single quotes (`'`) for note content to avoid shell interpretation issues
2. When editing notes, prefer the `-n` flag for overwriting content or omit if you pipe in with stdin
3. Use the list command to find note IDs for editing or deleting
4. Pipe content from files for longer notes or when special characters are needed
