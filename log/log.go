package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

const (
	Ldate         = 1 << iota                                 // the date in the local time zone: 2009/01/23
	Ltime                                                     // the time in the local time zone: 01:23:23
	Lmicroseconds                                             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                                 // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                                // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                                      // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                                                // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = log.LstdFlags | log.Lmsgprefix | log.LUTC // initial values for the standard logger
)

type Level int

const (
	ERROR Level = iota + 1
	WARNING
	INFO
)

var bold = color.New(color.Bold)

var tags = [4]string{
	bold.Add(color.FgRed).Sprint("[ERR]"),
	bold.Add(color.FgYellow).Sprint("[WRN]"),
	bold.Add(color.FgGreen).Sprint("[INF]"),
}

type Logger struct {
	writer    io.Writer
	save      bool
	prefix    string
	level     Level
	maxLevel  Level
	removeTag atomic.Bool
	root      log.Logger
}

func CreateWriter(file string) (writer *os.File, err error) {
	writer, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return
}

func New(prefix string, flag int, removeTag bool, filePath string, maxLevel Level) *Logger {
	lgr := &Logger{
		writer:   os.Stderr,
		prefix:   prefix,
		level:    WARNING,
		maxLevel: maxLevel,
	}
	if removeTag {
		lgr.removeTag.Store(true)
	}
	if filePath != "" {
		writer, err := CreateWriter(filePath)
		if err != nil {
			lgr.root.Printf("Could not write log to file %s, got error: %v", filePath, err)
		} else {
			lgr.writer = writer
			lgr.save = true
		}
	}

	fullPrefix := lgr.getFullPrefix()
	lgr.root = *log.New(lgr.writer, fullPrefix, flag)

	log.SetPrefix(fullPrefix)
	log.SetFlags(flag)

	return lgr
}

func (lgr *Logger) getFullPrefix() string {
	if lgr.removeTag.Load() {
		return lgr.prefix
	}
	return tags[lgr.level-1] + " " + lgr.prefix
}

func (lgr *Logger) logWithLevel(level Level) *log.Logger {
	lgr.level = level
	fullPrefix := lgr.getFullPrefix()
	lgr.root.SetPrefix(fullPrefix)
	if lgr.level > lgr.maxLevel {
		lgr.root.SetOutput(io.Discard)
	}
	return &lgr.root
}

func (lgr *Logger) WithInfoLevel() *log.Logger {
	return lgr.logWithLevel(INFO)
}

func (lgr *Logger) WithWarningLevel() *log.Logger {
	return lgr.logWithLevel(WARNING)
}

func (lgr *Logger) WithErrorLevel() *log.Logger {
	return lgr.logWithLevel(ERROR)
}

func (lgr *Logger) print(level Level, isDebug bool, format string, v ...any) {
	if lgr.maxLevel < level && !isDebug {
		return
	}
	if lgr.level != level {
		lgr.level = level
		fullPrefix := lgr.getFullPrefix()
		lgr.root.SetPrefix(fullPrefix)
	}
	if !lgr.save && lgr.level < INFO {
		lgr.root.SetOutput(os.Stderr)
	} else if !lgr.save {
		lgr.root.SetOutput(os.Stdout)
	}
	lgr.root.Output(3, fmt.Sprintf(format, v...))
}

func (lgr *Logger) Info(format string, v ...any) {
	lgr.print(INFO, false, format, v...)
}

func (lgr *Logger) Warning(format string, v ...any) {
	lgr.print(WARNING, false, format, v...)
}

func (lgr *Logger) Error(format string, v ...any) {
	lgr.print(ERROR, false, format, v...)
}

func (lgr *Logger) Panic(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	lgr.print(ERROR, false, s)
	panic(s)
}

func (lgr *Logger) Debug(level Level, format string, v ...any) {
	lgr.print(level, true, format, v...)
}

var std = New("", LstdFlags, false, "", WARNING)

func Default(level string, file ...string) *Logger {
	std.level = ConvertLevelString(level)
	if len(file) < 1 {
		writer, err := CreateWriter(file[0])
		if err != nil {
			std.root.Printf("Could not write log to file %s, got error: %v", file[0], err)
		} else {
			std.writer = writer
			std.save = true
		}
	}
	return std
}

func Warning(format string, v ...any) {
	std.print(WARNING, false, format, v...)
}

func ConvertLevelString(s string) Level {
	switch s {
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "INFO":
		return INFO
	default:
		panic(s + " is invalid value!")
	}
}
