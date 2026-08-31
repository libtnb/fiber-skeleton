// Command gen scaffolds a CRUD module (biz entity + repo interface, data repo
// implementation, service handlers, request structs, create-table migration
// and Wire module) or a standalone schema migration.
//
// Usage:
//
//	go run ./cmd/gen <name>            new module (singular snake_case, e.g. article)
//	go run ./cmd/gen migration <name>  new migration (e.g. add_email_to_users_table)
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

var (
	namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	// alterTarget extracts the table from names such as add_email_to_users.
	alterTarget = regexp.MustCompile(`_(?:to|from|in|on)_([a-z0-9_]+)$`)

	errUsage = errors.New(`usage:
  go run ./cmd/gen <name>            scaffold a CRUD module (singular snake_case, e.g. article)
  go run ./cmd/gen migration <name>  scaffold a schema migration (e.g. add_email_to_users_table)`)
)

type module struct {
	Module string // module path from go.mod
	Snake  string // article, order_item
	Pascal string // Article, OrderItem
	Camel  string // article, orderItem
	Table  string // articles, order_items
	Route  string // articles, order_items
	Date   string // 20260708120000, migration name prefix
}

// migration feeds migration.tmpl; the name doubles as the file name.
type migration struct {
	Name   string // 20260708120000_create_articles_table
	Table  string // articles
	Create bool   // create-table skeleton instead of an alter skeleton
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	switch {
	case len(args) == 2 && args[0] == "migration":
		return generateMigration(args[1])
	case len(args) == 1 && args[0] != "migration":
		return generateModule(args[0])
	default:
		return errUsage
	}
}

func generateModule(name string) error {
	if !namePattern.MatchString(name) {
		return errUsage
	}
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
	mig := migration{Name: m.Date + "_create_" + m.Table + "_table", Table: m.Table, Create: true}

	files := []struct {
		src  string
		dst  string
		data any
	}{
		{"biz.tmpl", filepath.Join("internal", m.Snake, "biz", m.Snake+".go"), m},
		{"data.tmpl", filepath.Join("internal", m.Snake, "data", m.Snake+".go"), m},
		{"service.tmpl", filepath.Join("internal", m.Snake, "service", "service.go"), m},
		{"request.tmpl", filepath.Join("internal", m.Snake, "service", "request.go"), m},
		{"route.tmpl", filepath.Join("internal", m.Snake, "service", "route.go"), m},
		{"module.tmpl", filepath.Join("internal", m.Snake, "wire.go"), m},
		{"migration.tmpl", filepath.Join("internal", "migrations", mig.Name+".go"), mig},
	}

	// refuse to overwrite anything: check all targets before writing any
	for _, f := range files {
		if _, err := os.Stat(f.dst); err == nil {
			return fmt.Errorf("%s already exists", f.dst)
		}
	}

	for _, f := range files {
		if err := render(f.src, f.dst, f.data); err != nil {
			return err
		}
		fmt.Println("created", f.dst)
	}

	fmt.Printf(`
Next steps:
  1. internal/app/wire.go: import "%[2]s/internal/%[1]s" and add
     "%[1]s.Module" to ApplicationModule.Include.
  2. run "make generate" to regenerate Wire code and mocks.
  3. run "make gen-check" to verify the generated module compiles.
     Mockery auto-discovers the new biz package and
     writes its repo mock under mocks/%[1]s/biz (no .mockery.yaml edit needed).
`, m.Snake, m.Module)

	return nil
}

func generateMigration(name string) error {
	if !namePattern.MatchString(name) {
		return errUsage
	}

	mig := migrationFor(time.Now().Format("20060102150405"), name)
	dst := filepath.Join("internal", "migrations", mig.Name+".go")
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%s already exists", dst)
	}
	if err := render("migration.tmpl", dst, mig); err != nil {
		return err
	}
	fmt.Println("created", dst)

	return nil
}

// migrationFor derives the skeleton from the migration name, following the
// create_<table>_table / add_<column>_to_<table>_table naming convention.
func migrationFor(date, name string) migration {
	mig := migration{Name: date + "_" + name, Table: "CHANGE_ME"}
	base := strings.TrimSuffix(name, "_table")
	if table, ok := strings.CutPrefix(base, "create_"); ok {
		mig.Table, mig.Create = table, true
	} else if sub := alterTarget.FindStringSubmatch(base); sub != nil {
		mig.Table = sub[1]
	}
	return mig
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

func render(src, dst string, data any) error {
	t, err := template.ParseFS(templates, "templates/"+src)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return err
	}
	code, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format %s: %w", dst, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600) //nolint:gosec // dst is namePattern-sanitized
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
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}

func toCamel(snake string) string {
	pascal := toPascal(snake)
	return strings.ToLower(pascal[:1]) + pascal[1:]
}
