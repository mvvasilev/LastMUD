package logging

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogLevel byte

const (
	DBG LogLevel = iota
	INF
	WRN
	ERR
)

func (ll LogLevel) String() string {
	switch ll {
	case DBG:
		return "DBG"
	case INF:
		return "INF"
	case WRN:
		return "WRN"
	case ERR:
		return "ERR"
	default:
		return strconv.FormatInt(int64(ll), 10)
	}
}

const DefaultLogFilePath = "./log/lastmud.log"
const DefaultTimeFormat = "2006-01-02 15:04:05.000"
const DefaultLogFileTimeFormat = "2006-01-02_15-04-05"

var DefaultLogger = CreateLogger(DBG, DBG, DefaultLogFilePath, DefaultTimeFormat, DefaultLogFileTimeFormat)

// Regex for color codes
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type Logger struct {
	file              *os.File
	timestampFormat   string
	maxFileLevel      LogLevel
	maxDisplayedLevel LogLevel
}

// Create a [Logger]
//
// Parameters:
//   - maxFileLevel: The maximum level to log ( in file )
//   - maxDisplayedLevel: The maximum level to display in the stdout ( can be different from [maxLevel] )
//   - filePath: The file to log to. Can be empty string if no file logging is desired.
func CreateLogger(maxFileLevel LogLevel, maxDisplayedLevel LogLevel, filePath string, timestampFormat string, fileTimestampFormat string) (logger *Logger) {
	logger = &Logger{}

	if filePath != "" {
		timestamp := time.Now().Format(fileTimestampFormat)

		ext := filepath.Ext(filePath)             // ".txt"
		base := strings.TrimSuffix(filePath, ext) // "./base/dir/log"

		logFilePath := fmt.Sprintf("%s-%s%s", base, timestamp, ext) // "./base/dir/log-2006-01-02_15-04-05.txt"

		mkdirErr := os.MkdirAll(filepath.Dir(logFilePath), 0755)

		if mkdirErr != nil {
			err(os.Stdout, false, timestampFormat, mkdirErr)
		}

		file, fileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if fileErr != nil {
			err(os.Stdout, false, timestampFormat, fileErr)
		} else {
			logger.file = file
		}
	}

	logger.timestampFormat = timestampFormat
	logger.maxFileLevel = maxFileLevel
	logger.maxDisplayedLevel = maxDisplayedLevel

	return
}

func (l *Logger) Debug(v ...any) {
	l.Log(DBG, v...)
}

func (l *Logger) Info(v ...any) {
	l.Log(INF, v...)
}

func (l *Logger) Warn(v ...any) {
	l.Log(WRN, v...)
}

func (l *Logger) Error(v ...any) {
	l.Log(ERR, v...)
}

func (l *Logger) Log(level LogLevel, v ...any) {
	if level >= l.maxDisplayedLevel {
		log(os.Stdout, level, false, l.timestampFormat, v...)
	}

	if level >= l.maxFileLevel {
		log(l.file, level, true, l.timestampFormat, v...)
	}
}

func Debug(v ...any) {
	DefaultLogger.Log(DBG, v...)
}

func Info(v ...any) {
	DefaultLogger.Log(INF, v...)
}

func Warn(v ...any) {
	DefaultLogger.Log(WRN, v...)
}

func Error(v ...any) {
	DefaultLogger.Log(ERR, v...)
}

func Log(level LogLevel, v ...any) {
	DefaultLogger.Log(level, v...)
}

func dbg(file *os.File, stripColor bool, timestampFormat string, v ...any) {
	custom(file, ColorWhite+DBG.String()+ColorReset, stripColor, timestampFormat, v...)
}

func inf(file *os.File, stripColor bool, timestampFormat string, v ...any) {
	custom(file, ColorCyan+INF.String()+ColorReset, stripColor, timestampFormat, v...)
}

func wrn(file *os.File, stripColor bool, timestampFormat string, v ...any) {
	custom(file, ColorYellow+WRN.String()+ColorReset, stripColor, timestampFormat, v...)
}

func err(file *os.File, stripColor bool, timestampFormat string, v ...any) {
	custom(file, ColorRed+ERR.String()+ColorReset, stripColor, timestampFormat, v...)
}

func custom(file *os.File, prefix string, stripColor bool, timestampFormat string, v ...any) {
	w := bufio.NewWriter(file)

	msg := fmt.Sprint(v...)

	if stripColor {
		msg = ansiRegex.ReplaceAllString(msg, "")
		prefix = ansiRegex.ReplaceAllString(prefix, "")
	}

	fmt.Fprintf(w, "%s | %s | %s\r\n", time.Now().Format(timestampFormat), prefix, msg)

	w.Flush()
}

func log(file *os.File, level LogLevel, stripColor bool, timestampFormat string, v ...any) {
	switch level {
	case DBG:
		dbg(file, stripColor, timestampFormat, v...)
	case INF:
		inf(file, stripColor, timestampFormat, v...)
	case WRN:
		wrn(file, stripColor, timestampFormat, v...)
	case ERR:
		err(file, stripColor, timestampFormat, v...)
	default:
		custom(file, level.String(), stripColor, timestampFormat, v...)
	}
}
