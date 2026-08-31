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

// nonModulePackages are the fixed top-level packages under internal/;
// every other directory is a business module, discovered automatically.
var nonModulePackages = map[string]bool{
	"app":        true,
	"migrations": true,
	"platform":   true,
	"shared":     true,
}

// TestModuleBoundaries fails when an internal import crosses the layering
// rules encoded in violation.
func TestModuleBoundaries(t *testing.T) {
	internalDir := filepath.Join("..", "..", "internal")

	entries, err := os.ReadDir(internalDir)
	if err != nil {
		t.Fatalf("read internal/: %v", err)
	}
	modules := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() && !nonModulePackages[e.Name()] {
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
		if len(segs) > 2 { // internal/<top>/<sub>/file.go
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
	targetSub := ""
	if len(targetSegs) > 1 {
		targetSub = targetSegs[1]
	}

	switch {
	case ownerTop == "app":
		return "" // the assembly imports everything

	case ownerTop == "migrations":
		return "migrations declares schema only and imports no internal package"

	case ownerTop == "shared":
		if targetTop == "shared" {
			return ""
		}
		return "shared holds the bottom-layer contracts and cannot depend on layers above"

	case ownerTop == "platform":
		switch {
		case ownerSub == "conf":
			return "platform/conf is the bottom layer and imports no internal package"
		case targetTop == "shared":
			return ""
		case targetTop == "platform" && (targetSub == ownerSub || targetSub == "conf"):
			return "" // own subtree, or the configuration everything reads
		case ownerSub == "bootstrap" && targetTop == "migrations":
			return "" // bootstrap hands the migration collection to the migrator
		}
		return "platform assembles infrastructure and must not know business modules"

	case modules[ownerTop]:
		switch {
		case targetTop == ownerTop:
			if ownerSub == "biz" && targetSub != "biz" {
				return "biz is the core and cannot import its own data/service adapters"
			}
			return ""
		case targetTop == "shared":
			return ""
		case modules[targetTop]:
			if targetSub == "biz" {
				return ""
			}
			return "modules reach each other only through the other module's biz package"
		default:
			return "modules depend on shared contracts and other modules' biz packages only"
		}
	}
	return ""
}
