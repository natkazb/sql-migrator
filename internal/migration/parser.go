package migration

import (
	"fmt"
	"strings"
)

const (
	FormatSQL     = "SQL"
	FormatGO      = "GO"
	FormatUnknown = "Unknown"
	upStart       = "-- Up begin"
	upEnd         = "-- Up end"
	downStart     = "-- Down begin"
	downEnd       = "-- Down end"
	sqlStart      = "-- SQL"
	goStart       = "-- GO"
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
		FormatUnknown: func(_ string) (*Migration, error) {
			return &Migration{
				Format: FormatUnknown,
			}, fmt.Errorf("unknown migration format")
		},
	}
)

func (p *BaseParser) Parse(content string) (*Migration, error) {
	currFormat := FormatUnknown
	if strings.Contains(content, sqlStart) {
		currFormat = FormatSQL
	} else if strings.Contains(content, goStart) {
		currFormat = FormatGO
	}
	return Parsers[currFormat](content)
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
		Format: "Base",
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
