package migration

import (
	"fmt"
	"strings"
)

const (
	template = `-- SQL
-- Up begin
-- SQL statements for applying the migration
-- Up end

-- Down begin
-- SQL statements for rolling back the migration
-- Down end
`
	templateGo = `-- GO
-- Up begin
-- GO statements
-- Up end

-- Down begin
-- GO statements
-- Down end
`
	FormatSQL = "SQL"
	FormatGO  = "GO"
	upStart   = "-- Up begin"
	upEnd     = "-- Up end"
	downStart = "-- Down begin"
	downEnd   = "-- Down end"
	sqlStart  = "-- SQL"
	goStart   = "-- GO"
)

type Migration struct {
	Format string
	Up     string
	Down   string
}

func ParseMigration(content string) (*Migration, error) {
	upStartIdx := strings.Index(content, upStart)
	upEndIdx := strings.Index(content, upEnd)
	downStartIdx := strings.Index(content, downStart)
	downEndIdx := strings.Index(content, downEnd)

	if upStartIdx == -1 || upEndIdx == -1 || downStartIdx == -1 || downEndIdx == -1 {
		return nil, fmt.Errorf("invalid migration format")
	}

	upSQL := strings.TrimSpace(content[upStartIdx+len(upStart) : upEndIdx])
	downSQL := strings.TrimSpace(content[downStartIdx+len(downStart) : downEndIdx])

	return &Migration{
		Up:   upSQL,
		Down: downSQL,
	}, nil
}

type Migrate interface {
	Parse(_ string)
}

type GoMigrate struct{}

func (m *GoMigrate) Parse(_ string) {
}

type SQLMigrate struct{}

func (m *SQLMigrate) Parse(_ string) {
}
