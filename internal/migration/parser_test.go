package migration_test

import (
	"os"
	"testing"

	"github.com/natkazb/sql-migrator/internal/migration" //nolint:depguard
	"github.com/stretchr/testify/assert"
)

func TestBaseParser_Parse(t *testing.T) {
	parser := &migration.BaseParser{}

	t.Run("Valid SQL File", func(t *testing.T) {
		filePath := "test.sql"
		content := `-- Up begin
CREATE TABLE users (id SERIAL PRIMARY KEY);
-- Up end
-- Down begin
DROP TABLE users;
-- Down end`
		err := os.WriteFile(filePath, []byte(content), 0o644)
		assert.NoError(t, err)
		defer os.Remove(filePath)

		mig, err := parser.Parse(filePath)
		assert.NoError(t, err)
		assert.Equal(t, migration.FormatSQL, mig.Format)
		assert.Equal(t, "CREATE TABLE users (id SERIAL PRIMARY KEY);", mig.Up)
		assert.Equal(t, "DROP TABLE users;", mig.Down)
	})

	t.Run("Unknown File Format", func(t *testing.T) {
		filePath := "test.unknown"
		content := `-- Unknown format`
		err := os.WriteFile(filePath, []byte(content), 0o644)
		assert.NoError(t, err)
		defer os.Remove(filePath)

		mig, err := parser.Parse(filePath)
		assert.Error(t, err)
		assert.Equal(t, migration.FormatUnknown, mig.Format)
	})
}

func TestSQLParser_Parse(t *testing.T) {
	parser := &migration.SQLParser{}

	t.Run("Valid SQL Migration", func(t *testing.T) {
		content := `-- Up begin
CREATE TABLE users (id SERIAL PRIMARY KEY);
-- Up end
-- Down begin
DROP TABLE users;
-- Down end`
		mig, err := parser.Parse(content)
		assert.NoError(t, err)
		assert.Equal(t, "CREATE TABLE users (id SERIAL PRIMARY KEY);", mig.Up)
		assert.Equal(t, "DROP TABLE users;", mig.Down)
	})

	t.Run("Invalid SQL Migration", func(t *testing.T) {
		content := `-- Up begin
CREATE TABLE users (id SERIAL PRIMARY KEY);
-- Down begin
DROP TABLE users;`
		mig, err := parser.Parse(content)
		assert.Error(t, err)
		assert.Nil(t, mig)
	})
}

func TestGOParser_Parse(t *testing.T) {
	parser := &migration.GOParser{}

	t.Run("Not Implemented", func(t *testing.T) {
		content := `-- Up begin
-- GO statements
-- Up end

-- Down begin
-- GO statements
-- Down end`
		mig, err := parser.Parse(content)
		assert.Error(t, err)
		assert.Equal(t, migration.FormatGO, mig.Format)
	})
}
