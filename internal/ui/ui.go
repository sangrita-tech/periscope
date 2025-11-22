package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	Title   = color.New(color.FgCyan, color.Bold)
	Info    = color.New(color.FgHiBlue)
	Success = color.New(color.FgHiGreen)
	Warn    = color.New(color.FgHiYellow)
	Error   = color.New(color.FgHiRed, color.Bold)

	Dim = color.New(color.Faint)

	DirStyle  = color.New(color.FgHiCyan, color.Bold)
	FileStyle = color.New(color.FgWhite)
	Branch    = color.New(color.FgHiBlack)

	FileHeader = color.New(color.FgMagenta, color.Bold)
)

type Logger struct {
	Out io.Writer
	Err io.Writer
}

func New() *Logger {
	return &Logger{
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

func (l *Logger) PrintTitle(format string, a ...any) {
	_, _ = fmt.Fprintln(l.Out, Title.Sprintf("🔭 "+format, a...))
}
func (l *Logger) Info(format string, a ...any) {
	_, _ = fmt.Fprintln(l.Out, Info.Sprintf("ℹ️  "+format, a...))
}
func (l *Logger) Success(format string, a ...any) {
	_, _ = fmt.Fprintln(l.Out, Success.Sprintf("✅ "+format, a...))
}
func (l *Logger) Warn(format string, a ...any) {
	_, _ = fmt.Fprintln(l.Err, Warn.Sprintf("⚠️  "+format, a...))
}
func (l *Logger) Error(format string, a ...any) {
	_, _ = fmt.Fprintln(l.Err, Error.Sprintf("❌ "+format, a...))
}

func (l *Logger) PlainErr(err error) {
	_, _ = fmt.Fprintln(l.Err, err.Error())
}
