package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	return &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func Info(format string, v ...any) {
	color.Cyan(format, v...)
}

func Warn(format string, v ...any) {
	color.Yellow(format, v...)
}

func Error(err error) {
	color.Red(err.Error())
}

func Link(format string, v ...any) {
	link := fmt.Sprintf(format, v...)
	color.Blue(termlink.Link(link, link))
}
