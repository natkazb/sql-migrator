package migration

type Migrate interface {
	Parse(_ string)
}

type GoMigrate struct{}

func (m *GoMigrate) Parse(_ string) {
}

type SQLMigrate struct{}

func (m *SQLMigrate) Parse(_ string) {
}
