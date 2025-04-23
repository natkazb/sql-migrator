package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/natkazb/sql-migrator/internal/config"    //nolint:depguard
)

func TestNewConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	tempFile, err := os.CreateTemp("", "config_test_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write valid YAML content to the temporary file
	yamlContent := `
logger:
  level: "info"
sql:
  host: "localhost"
  port: 5432
  dbName: "testdb"
  user: "testuser"
  password: "testpassword"
  driver: "postgres"
path: "/some/path"
`
	_, err = tempFile.WriteString(yamlContent)
	assert.NoError(t, err)
	tempFile.Close()

	// Test the NewConfig function with the valid YAML file
	conf, err := config.NewConfig(tempFile.Name())
	assert.NoError(t, err)

	// Verify the returned Config struct
	expectedConfig := config.Config{
		Logger: config.LoggerConf{Level: "info"},
		SQL: config.SQLConf{
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
			Username: "testuser",
			Password: "testpassword",
			Driver:   "postgres",
		},
		Path: "/some/path",
	}
	assert.Equal(t, expectedConfig, conf)

	// Test error case: file does not exist
	_, err = config.NewConfig("non_existent_file.yaml")
	assert.Error(t, err)

	// Test error case: invalid YAML content
	invalidTempFile, err := os.CreateTemp("", "invalid_config_test_*.yaml")
	assert.NoError(t, err)
	defer os.Remove(invalidTempFile.Name())

	_, err = invalidTempFile.WriteString("invalid_yaml: [")
	assert.NoError(t, err)
	invalidTempFile.Close()

	_, err = config.NewConfig(invalidTempFile.Name())
	assert.Error(t, err)
}