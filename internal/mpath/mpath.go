package mpath

import (
	"fmt"
	"os"
	"sort"
	"time"
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
