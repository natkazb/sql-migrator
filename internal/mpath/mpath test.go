package mpath

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	Messages []string
}

func (l *MockLogger) Info(msg string)  { l.Messages = append(l.Messages, "[INFO]: "+msg) }
func (l *MockLogger) Debug(msg string) { l.Messages = append(l.Messages, "[DEBUG]: "+msg) }
func (l *MockLogger) Warn(msg string)  { l.Messages = append(l.Messages, "[WARN]: "+msg) }
func (l *MockLogger) Error(msg string) { l.Messages = append(l.Messages, "[ERROR]: "+msg) }

func TestCreateNew(t *testing.T) {
	tempDir := t.TempDir()
	logger := &MockLogger{}
	migrationPath := New(tempDir, logger)

	migrationName := "test_migration"
	migrationPath.CreateNew(migrationName, "SQL")

	files, err := os.ReadDir(tempDir)
	require.NoError(t, err, "Failed to read temp directory: %v", tempDir)

	require.Equal(t, 1, len(files), "Expected 1 file, got %d", len(files))

	expectedFileName := time.Now().Format(formatTime) + "_" + migrationName + ".sql"
	require.Equal(t, expectedFileName, files[0].Name(), "Expected file name %s, got %s", expectedFileName, files[0].Name())

	content, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
	require.NoError(t, err, "Failed to read file content: %v", err)

	require.Equal(t, template, string(content), "File content does not match template. Got: %s", string(content))
}

func TestGetList(t *testing.T) {
	tempDir := t.TempDir()
	logger := &MockLogger{}
	migrationPath := New(tempDir, logger)

	// Create test files with different modification times
	file1 := filepath.Join(tempDir, "file1.sql")
	file2 := filepath.Join(tempDir, "file2.sql")
	file3 := filepath.Join(tempDir, "file3.sql")

	os.WriteFile(file1, []byte{}, fs.ModePerm)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(file2, []byte{}, fs.ModePerm)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(file3, []byte{}, fs.ModePerm)

	files, err := migrationPath.GetList()
	require.NoError(t, err, "GetList failed: %v", err)

	expectedOrder := []string{"file1.sql", "file2.sql", "file3.sql"}
	for i, file := range files {
		require.Equal(t, expectedOrder[i], file, "Expected %s, got %s", expectedOrder[i], file)
	}
}
