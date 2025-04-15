package migrator

type Migrate interface {
    Parse(message string)
}

type GoMigrate struct{}

func (m *GoMigrate) Parse(message string) {
}

type SqlMigrate struct{}

func (m *SqlMigrate) Parse(message string) {
}