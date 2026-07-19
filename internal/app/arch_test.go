package app_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// modulePrefix is this repository's import path for internal packages.
const modulePrefix = "github.com/libtnb/fiber-skeleton/internal/"

// sharedPackages are the non-module top-level packages under internal/;
// everything else is a business module, discovered automatically.
var sharedPackages = map[string]bool{
	"app":       true,
	"bootstrap": true,
	"conf":      true,
	"pkg":       true,
	"server":    true,
}

// TestModuleBoundaries fails CI when an import crosses the monolith's rules:
// modules reach each other only through biz, never import the composition
// layers or conf, internal/pkg depends on nothing above it, and biz stays
// free of its own adapters.
func TestModuleBoundaries(t *testing.T) {
	internalDir := filepath.Join("..", "..", "internal")

	entries, err := os.ReadDir(internalDir)
	if err != nil {
		t.Fatalf("read internal/: %v", err)
	}
	modules := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() && !sharedPackages[e.Name()] {
			modules[e.Name()] = true
		}
	}
	if len(modules) == 0 {
		t.Fatal("no business modules discovered under internal/")
	}

	fset := token.NewFileSet()
	err = filepath.WalkDir(internalDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}

		rel, err := filepath.Rel(internalDir, path)
		if err != nil {
			return err
		}
		segs := strings.Split(filepath.ToSlash(rel), "/")
		ownerTop := segs[0]
		ownerSub := ""
		if len(segs) > 2 { // internal/<module>/<layer>/file.go
			ownerSub = segs[1]
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range f.Imports {
			target := strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(target, modulePrefix) {
				continue
			}
			if msg := violation(modules, ownerTop, ownerSub, strings.TrimPrefix(target, modulePrefix)); msg != "" {
				t.Errorf("%s imports %s: %s", rel, target, msg)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk internal/: %v", err)
	}
}

// violation returns why this import edge is illegal, or "".
func violation(modules map[string]bool, ownerTop, ownerSub, target string) string {
	targetSegs := strings.Split(target, "/")
	targetTop := targetSegs[0]

	switch {
	case ownerTop == "app":
		return "" // the assembly imports everything

	case ownerTop == "conf":
		return "conf is the bottom layer and imports no internal package"

	case ownerTop == "pkg":
		if targetTop != "pkg" {
			return "internal/pkg holds shared contracts and cannot depend on layers above it"
		}
		return ""

	case ownerTop == "bootstrap" || ownerTop == "server":
		if targetTop == "pkg" || targetTop == "conf" {
			return ""
		}
		return ownerTop + " assembles infrastructure and must not know business modules"

	case modules[ownerTop]:
		switch {
		case targetTop == ownerTop:
			if ownerSub == "biz" && len(targetSegs) > 1 && targetSegs[1] != "biz" {
				return "biz is the core and cannot import its own data/service adapters"
			}
			return ""
		case targetTop == "pkg":
			return ""
		case modules[targetTop]:
			if len(targetSegs) > 1 && targetSegs[1] == "biz" {
				return ""
			}
			return "modules reach each other only through the other module's biz package"
		default:
			return "modules depend on internal/pkg and other modules' biz packages only"
		}
	}
	return ""
}
