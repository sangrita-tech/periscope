package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	Title      = color.New(color.FgHiMagenta, color.Bold)
	Info       = color.New(color.FgMagenta)
	Success    = color.New(color.FgHiCyan, color.Bold)
	Warn       = color.New(color.FgHiMagenta)
	Error      = color.New(color.FgMagenta, color.Bold)
	Dim        = color.New(color.FgHiBlack)
	DirStyle   = color.New(color.FgHiMagenta, color.Bold)
	FileStyle  = color.New(color.FgHiWhite)
	Branch     = color.New(color.FgHiBlack)
	FileHeader = color.New(color.FgHiMagenta, color.Bold)
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
