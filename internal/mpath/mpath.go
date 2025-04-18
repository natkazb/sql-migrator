package mpath

import (
	"os"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type MigrationPath struct {
	Path string
	log  Logger
}

type MigrationFiles []string

func New(path string, l Logger) *MigrationPath {
	return &MigrationPath{
		Path: path,
		log:  l,
	}
}

func (m *MigrationPath) GetList() ([]string, error) {
	files := make([]string, 0) // @todo: можно предварительно посчитать кол-во файлов в директории

	entries, err := os.ReadDir(m.Path)
	if err != nil {
		m.log.Error("Failed to read directory: " + err.Error())
		return files, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}
