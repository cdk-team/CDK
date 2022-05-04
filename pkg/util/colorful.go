package util

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

type Level uint8

const (
	ERROR Level = iota
	WARNNING
	INFO
	DEBUG
)

var DefaultLevel = INFO

const (
	DebugPrefix = "[DEBUG] "
	InfoPrefix  = "[INFO_]  "
	WarnPrefix  = "[WARN_]  "
	ErrorPrefix = "[ERROR] "
)

// Colorful Bold
// use like `GreenBold.Sprint(str)`
var (
	RedBold    = color.New(color.FgRed).Add(color.Bold)
	GreenBold  = color.New(color.FgGreen).Add(color.Bold)
	YellowBold = color.New(color.FgYellow).Add(color.Bold)
	BlueBold   = color.New(color.FgBlue).Add(color.Bold)
)

type LevelLogger struct {
	Level Level
	Color bool

	PrintFunc func(format string, v ...interface{})
}

var (
	ColorDebugPrefix = GreenBold.Sprint(DebugPrefix)
	ColorInfoPrefix  = BlueBold.Sprint(InfoPrefix)
	ColorWarnPrefix  = YellowBold.Sprint(WarnPrefix)
	ColorErrorPrefix = RedBold.Sprint(ErrorPrefix)
)

func (l *LevelLogger) Debug(format string, v ...interface{}) {
	if l.Level >= DEBUG {
		prefix := DebugPrefix
		if l.Color {
			prefix = ColorDebugPrefix
		}

		temp := fmt.Sprintf("%s%s", prefix, format)
		l.PrintFunc(temp, v...)
	}
}

func (l *LevelLogger) Info(format string, v ...interface{}) {
	if l.Level >= INFO {
		prefix := InfoPrefix
		if l.Color {
			prefix = ColorInfoPrefix
		}

		temp := fmt.Sprintf("%s%s", prefix, format)
		l.PrintFunc(temp, v...)
	}
}

func (l *LevelLogger) Warn(format string, v ...interface{}) {
	if l.Level >= WARNNING {
		prefix := WarnPrefix
		if l.Color {
			prefix = ColorWarnPrefix
		}

		temp := fmt.Sprintf("%s%s", prefix, format)
		l.PrintFunc(temp, v...)
	}
}

func (l *LevelLogger) Error(format string, v ...interface{}) {
	prefix := ErrorPrefix
	if l.Color {
		prefix = ColorErrorPrefix
	}

	temp := fmt.Sprintf("%s%s", prefix, format)
	l.PrintFunc(temp, v...)
}

func (l *LevelLogger) Close() {}

type Wrapper struct {
	logger *log.Logger

	LevelLogger
}

func NewWrapper(writer io.Writer, colorful bool) *Wrapper {
	logger := log.New(writer, "", log.LstdFlags|log.Lshortfile)

	return &Wrapper{
		logger: logger,
		LevelLogger: LevelLogger{
			Level:     DefaultLevel,
			Color:     colorful,
			PrintFunc: logger.Printf,
		},
	}
}

func NewStdoutWrapper() *Wrapper {
	return NewWrapper(os.Stdout, true)
}
