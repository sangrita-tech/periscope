package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

// ViewOptions описывает всё, что нужно для выполнения команды view.
type ViewOptions struct {
	Dir             string
	CopyToClipboard bool
	StripComments   bool
	IgnorePath      []string
	IgnoreContent   []string
}

// RunView выполняет основную логику команды view.
func RunView(opts ViewOptions) error {
	info, err := os.Stat(opts.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q does not exist", opts.Dir)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied for directory %q", opts.Dir)
		}
		return fmt.Errorf("failed to access directory %q: %v", opts.Dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", opts.Dir)
	}

	var out strings.Builder

	// Куда пишем вывод (в буфер и/или в stdout)
	write := func(s string) {
		if opts.CopyToClipboard {
			out.WriteString(s)
		}
		if !opts.CopyToClipboard {
			fmt.Print(s)
		}
	}

	// Лог ошибок — всегда пишет в stderr, а при необходимости дублирует в буфер
	logErr := func(action, filePath string, err error) {
		abs, _ := filepath.Abs(filePath)
		msg := fmt.Sprintf("%s - %s (%v)\n", action, abs, err)
		if opts.CopyToClipboard {
			out.WriteString(msg)
		}
		fmt.Fprint(os.Stderr, msg)
	}

	firstFile := true

	walkErr := filepath.WalkDir(opts.Dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			logErr("access error", path, walkErr)
			return nil
		}

		// Сначала проверяем, надо ли игнорировать по пути (и файлы, и директории)
		if pathMatchesAny(opts.IgnorePath, path) {
			if d.IsDir() {
				// не заходить внутрь этой директории
				return fs.SkipDir
			}
			// просто пропускаем файл
			return nil
		}

		// дальше нас интересуют только файлы
		if d.IsDir() {
			return nil
		}

		if err := processFile(path, opts, &out, &firstFile, write, logErr); err != nil {
			logErr("processing failed", path, err)
			// не прерываем обход полностью
			return nil
		}

		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	if opts.CopyToClipboard {
		if err := clipboard.WriteAll(out.String()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %v", err)
		}
	}

	return nil
}

// processFile читает файл, применяет фильтры по содержимому и выводит его при необходимости.
func processFile(
	path string,
	opts ViewOptions,
	out io.StringWriter,
	firstFile *bool,
	write func(string),
	logErr func(action, filePath string, err error),
) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path failed: %w", err)
	}

	scanner := bufio.NewScanner(f)

	var fileBuf strings.Builder
	skipFile := false

	for scanner.Scan() {
		line := scanner.Text()

		// если содержимое попадает под ignoreContent — выкидываем файл целиком
		if lineMatchesAny(opts.IgnoreContent, line) {
			skipFile = true
			break
		}

		if opts.StripComments {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "--") {
				continue
			}
		}

		fileBuf.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	if skipFile {
		return nil
	}

	// файл прошёл все фильтры — выводим
	if *firstFile {
		write(fmt.Sprintf("[FILE] %s\n\n", absPath))
		*firstFile = false
	} else {
		write(fmt.Sprintf("\n[FILE] %s\n\n", absPath))
	}

	write(fileBuf.String())

	return nil
}

// pathMatchesAny проверяет путь и basename по всем паттернам.
func pathMatchesAny(patterns []string, path string) bool {
	if len(patterns) == 0 {
		return false
	}

	base := filepath.Base(path)

	for _, pat := range patterns {
		if pat == "" {
			continue
		}

		if matchPattern(pat, path) {
			return true
		}
		// отдельно проверяем basename, чтобы ".git" матчило директорию ".git"
		if base != path && matchPattern(pat, base) {
			return true
		}
	}

	return false
}

// lineMatchesAny проверяет строку по всем паттернам.
func lineMatchesAny(patterns []string, line string) bool {
	if len(patterns) == 0 {
		return false
	}

	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		if matchPattern(pat, line) {
			return true
		}
	}
	return false
}

// matchPattern — обёртка над filepath.Match: при ошибке паттерна просто не матчим.
func matchPattern(pattern, s string) bool {
	ok, err := filepath.Match(pattern, s)
	return err == nil && ok
}
