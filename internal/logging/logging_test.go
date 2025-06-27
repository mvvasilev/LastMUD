package logging

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelString(t *testing.T) {
	assert.Equal(t, "DBG", DBG.String())
	assert.Equal(t, "INF", INF.String())
	assert.Equal(t, "WRN", WRN.String())
	assert.Equal(t, "ERR", ERR.String())
	assert.Equal(t, "99", LogLevel(99).String())
}

func TestCreateLogger_NoFile(t *testing.T) {
	logger := CreateLogger(DBG, INF, "", "2006-01-02 15:04:05.000", "2006-01-02_15-04-05")
	assert.NotNil(t, logger)
	assert.Nil(t, logger.file)
	assert.Equal(t, DBG, logger.maxFileLevel)
	assert.Equal(t, INF, logger.maxDisplayedLevel)
}

func TestCreateLogger_WithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	logger := CreateLogger(DBG, DBG, path, "2006-01-02 15:04:05.000", "2006-01-02_15-04-05")
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.file)
}

func TestLoggerMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := &Logger{
		file:              nil,
		timestampFormat:   DefaultTimeFormat,
		maxFileLevel:      DBG,
		maxDisplayedLevel: DBG,
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

func TestAnsiStripRegex(t *testing.T) {
	input := "\033[31mHello\033[0m"
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	output := re.ReplaceAllString(input, "")
	assert.Equal(t, "Hello", output)
}
