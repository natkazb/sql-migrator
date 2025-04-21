package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		level       string
		method      func(*Logger, string)
		message     string
		expectedLog bool
	}{
		{"DEBUG", (*Logger).Debug, "Debug message", true},
		{"DEBUG", (*Logger).Info, "Info message", true},
		{"DEBUG", (*Logger).Warn, "Warn message", true},
		{"DEBUG", (*Logger).Error, "Error message", true},
		{"INFO", (*Logger).Debug, "Debug message", false},
		{"INFO", (*Logger).Info, "Info message", true},
		{"INFO", (*Logger).Warn, "Warn message", true},
		{"INFO", (*Logger).Error, "Error message", true},
		{"WARN", (*Logger).Debug, "Debug message", false},
		{"WARN", (*Logger).Info, "Info message", false},
		{"WARN", (*Logger).Warn, "Warn message", true},
		{"WARN", (*Logger).Error, "Error message", true},
		{"ERROR", (*Logger).Debug, "Debug message", false},
		{"ERROR", (*Logger).Info, "Info message", false},
		{"ERROR", (*Logger).Warn, "Warn message", false},
		{"ERROR", (*Logger).Error, "Error message", true},
	}

	for _, test := range tests {
		t.Run(test.level+"_"+test.message, func(t *testing.T) {
			// Сохранить оригинальный os.Stdout
			originalStdout := os.Stdout
			// Создать pipe
			r, w, err := os.Pipe()
			require.NoError(t, err, "Ошибка создания pipe: %v\n", err)
			if err != nil {
				return
			}
			// Перенаправить os.Stdout на запись в pipe
			os.Stdout = w

			// Create logger and call the method
			logger := New(test.level)
			test.method(logger, test.message)

			// Закрыть запись в pipe
			w.Close()

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			require.NoError(t, err, "Ошибка чтения из pipe: %v\n", err)
			if err != nil {
				return
			}

			// Восстановить оригинальный os.Stdout
			os.Stdout = originalStdout

			// Check if the message was logged
			output := buf.String()
			if test.expectedLog {
				require.Contains(t, output, test.message, "Expected log message '%s', but it was not logged", test.message)
			}
			if !test.expectedLog {
				require.NotContains(t, output, test.message, "Did not expect log message '%s', but it was logged", test.message)
			}
		})
	}
}
