package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	FormatSQL     = "sql"
	FormatGO      = "go"
	FormatUnknown = "Unknown"
	upStart       = "-- Up begin"
	upEnd         = "-- Up end"
	downStart     = "-- Down begin"
	downEnd       = "-- Down end"
)

type Migration struct {
	Format string
	Up     string
	Down   string
}

type Parser interface {
	Parse(content string) (*Migration, error)
}

type BaseParser struct{}

var (
	sqlParser = &SQLParser{}
	goParser  = &GOParser{}
	Parsers   = map[string]func(content string) (*Migration, error){
		FormatSQL: func(content string) (*Migration, error) {
			return sqlParser.Parse(content)
		},
		FormatGO: func(content string) (*Migration, error) {
			return goParser.Parse(content)
		},
	}
)

func (p *BaseParser) Parse(filePath string) (*Migration, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return &Migration{}, fmt.Errorf("failed to read file %s: %s", filePath, err.Error())
	}
	fileExt := strings.ToLower(strings.Trim(filepath.Ext(filePath), "."))
	currParser, ok := Parsers[fileExt]
	if ok {
		return currParser(string(content))
	}
	return &Migration{Format: FormatUnknown}, fmt.Errorf("unknown migration format")
}

type SQLParser struct{}

func (p *SQLParser) Parse(content string) (*Migration, error) {
	upStartIdx := strings.Index(content, "-- Up begin")
	upEndIdx := strings.Index(content, "-- Up end")
	downStartIdx := strings.Index(content, "-- Down begin")
	downEndIdx := strings.Index(content, "-- Down end")

	if upStartIdx == -1 || upEndIdx == -1 || downStartIdx == -1 || downEndIdx == -1 {
		return nil, fmt.Errorf("invalid migration format")
	}

	upSQL := strings.TrimSpace(content[upStartIdx+len("-- Up begin") : upEndIdx])
	downSQL := strings.TrimSpace(content[downStartIdx+len("-- Down begin") : downEndIdx])

	return &Migration{
		Format: FormatSQL,
		Up:     upSQL,
		Down:   downSQL,
	}, nil
}

type GOParser struct{}

func (p *GOParser) Parse(_ string) (*Migration, error) {
	return &Migration{
		Format: FormatGO,
	}, fmt.Errorf("not implemented")
}
