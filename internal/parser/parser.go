package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses a .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Comment-only line
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		entry, err := parseLine(trimmed, lineNum)
		if err != nil {
			return nil, fmt.Errorf("%s line %d: %w", path, lineNum, err)
		}

		env.Entries = append(env.Entries, entry)
		env.Index[entry.Key] = &env.Entries[len(env.Entries)-1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning %s: %w", path, err)
	}

	return env, nil
}

func parseLine(line string, lineNum int) (Entry, error) {
	// Strip inline comment
	comment := ""
	if idx := strings.Index(line, " #"); idx != -1 {
		comment = strings.TrimSpace(line[idx+1:])
		line = strings.TrimSpace(line[:idx])
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return Entry{}, fmt.Errorf("invalid format, expected KEY=VALUE")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

	if key == "" {
		return Entry{}, fmt.Errorf("empty key")
	}

	return Entry{Key: key, Value: value, Comment: comment, Line: lineNum}, nil
}
