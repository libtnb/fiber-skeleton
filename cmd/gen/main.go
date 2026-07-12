// Command gen scaffolds a CRUD module: biz entity + repo interface, data
// repo implementation, service handlers, request structs and a migration.
//
// Usage: go run ./cmd/gen <name>   (name is singular snake_case, e.g. article)
package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/jinzhu/inflection"
)

//go:embed templates/*.tmpl
var templates embed.FS

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type module struct {
	Module string // module path from go.mod
	Snake  string // article, order_item
	Pascal string // Article, OrderItem
	Camel  string // article, orderItem
	Table  string // articles, order_items
	Route  string // articles, order_items
	Date   string // 20260708120000, migration name prefix
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 || !namePattern.MatchString(os.Args[1]) {
		return errors.New("usage: go run ./cmd/gen <name> (singular snake_case, e.g. article or order_item)")
	}

	name := os.Args[1]
	modPath, err := modulePath()
	if err != nil {
		return err
	}

	m := module{
		Module: modPath,
		Snake:  name,
		Pascal: toPascal(name),
		Camel:  toCamel(name),
		Table:  inflection.Plural(name),
		Route:  inflection.Plural(name),
		Date:   time.Now().Format("20060102150405"),
	}

	files := map[string]string{
		"biz.tmpl":       filepath.Join("internal", m.Snake, "biz", m.Snake+".go"),
		"data.tmpl":      filepath.Join("internal", m.Snake, "data", m.Snake+".go"),
		"migration.tmpl": filepath.Join("internal", m.Snake, "data", "migration.go"),
		"service.tmpl":   filepath.Join("internal", m.Snake, "service", "service.go"),
		"request.tmpl":   filepath.Join("internal", m.Snake, "service", "request.go"),
		"route.tmpl":     filepath.Join("internal", m.Snake, "service", "route.go"),
		"module.tmpl":    filepath.Join("internal", m.Snake, m.Snake+".go"),
	}

	// refuse to overwrite anything: check all targets before writing any
	for _, dst := range files {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("%s already exists", dst)
		}
	}

	for src, dst := range files {
		if err := render(src, dst, m); err != nil {
			return err
		}
		fmt.Println("created", dst)
	}

	fmt.Printf(`
Next steps:
  1. internal/app/injector.go: import "%[2]s/internal/%[1]s" and add
     "%[1]s.Package," to the business modules list.
  2. run "make generate" — mockery auto-discovers the new biz package and
     writes its repo mock under mocks/biz (no .mockery.yaml edit needed).
`, m.Snake, m.Module)

	return nil
}

func modulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("read go.mod (run from the project root): %w", err)
	}
	for line := range strings.Lines(string(data)) {
		if after, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			return strings.TrimSpace(after), nil
		}
	}
	return "", errors.New("cannot determine module path from go.mod")
}

func render(src, dst string, m module) error {
	t, err := template.ParseFS(templates, "templates/"+src)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, m); err != nil {
		return err
	}
	code, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format %s: %w", dst, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = f.Write(code); err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

func toPascal(snake string) string {
	var b strings.Builder
	for part := range strings.SplitSeq(snake, "_") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	return b.String()
}

func toCamel(snake string) string {
	pascal := toPascal(snake)
	return strings.ToLower(pascal[:1]) + pascal[1:]
}
