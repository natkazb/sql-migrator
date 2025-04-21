package mpath

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

const (
	formatTime = "20060102150405"
	template   = `-- SQL
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
	FormatGO = "GO"
)

type MigrationPath struct {
	Path string
	log  Logger
}

type FileInfo struct {
	Name    string
	ModTime time.Time
}

func New(path string, l Logger) *MigrationPath {
	return &MigrationPath{
		Path: path,
		log:  l,
	}
}

func (m *MigrationPath) CreateNew(name, format string) {
	timestamp := time.Now().Format(formatTime)
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filePath := filepath.Join(m.Path, filename)
	file, err := os.Create(filePath)
	if err != nil {
		m.log.Error(fmt.Sprintf("Error creating file: %s : %s", filePath, err.Error()))
	}
	defer file.Close()

	template := template
	if format == FormatGO {
		template = templateGo
	}
	_, err = file.WriteString(template)
	if err != nil {
		m.log.Error(fmt.Sprintf("Error writing file: %s : %s", filePath, err.Error()))
	}
	m.log.Info(fmt.Sprintf("Migration created in %s", filePath))
}

func (m *MigrationPath) GetList() ([]string, error) {
	files := make([]string, 0)
	entries, err := os.ReadDir(m.Path)
	if err != nil {
		m.log.Error("Failed to read directory: " + err.Error())
		return files, err
	}

	// Collect file info with modification times
	fileInfos := make([]FileInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:    entry.Name(),
			ModTime: info.ModTime(),
		})
	}

	// Sort files by modification time
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime.Before(fileInfos[j].ModTime)
	})

	for _, fileInfo := range fileInfos {
		files = append(files, fileInfo.Name)
	}

	return files, nil
}
