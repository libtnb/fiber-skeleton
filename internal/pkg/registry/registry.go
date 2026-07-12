// Package registry collects same-kind contributions (routes, commands,
// jobs) registered under a naming convention, do's stand-in for value groups.
// Each module registers its routes/commands/jobs as named providers under one
// of the prefixes below; the server, CLI and scheduler collect them without
// knowing which module they came from.
package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samber/do/v2"
)

// Contribution prefixes; a module names its providers with these so the
// assemblers can gather them.
const (
	RoutePrefix      = "routes:"
	CommandPrefix    = "commands:"
	JobPrefix        = "jobs:"
	SubscriberPrefix = "subscribers:"
)

// Collect resolves every service whose name starts with prefix, sorted by
// name for deterministic order.
func Collect[T any](i do.Injector, prefix string) ([]T, error) {
	var names []string
	for _, desc := range i.ListProvidedServices() {
		if strings.HasPrefix(desc.Service, prefix) {
			names = append(names, desc.Service)
		}
	}
	sort.Strings(names)

	out := make([]T, 0, len(names))
	for _, name := range names {
		svc, err := do.InvokeNamed[T](i, name)
		if err != nil {
			return nil, err
		}
		out = append(out, svc)
	}

	return out, nil
}

// Verify fails on named contributions whose prefix matches none of the known
// ones — a typo like "route:user" would otherwise be silently dropped.
// Unnamed services (type-derived names, no colon) are ignored.
func Verify(i do.Injector, prefixes ...string) error {
	for _, desc := range i.ListProvidedServices() {
		name := desc.Service
		if !strings.Contains(name, ":") {
			continue
		}
		known := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				known = true
				break
			}
		}
		if !known {
			return fmt.Errorf("contribution %q matches no known prefix %v", name, prefixes)
		}
	}

	return nil
}

// Lazy adapts a plain single-dependency constructor into a lazy provider
// entry, keeping constructors container-free.
func Lazy[T, D any](ctor func(D) T) func(do.Injector) {
	return do.Lazy(func(i do.Injector) (T, error) {
		return ctor(do.MustInvoke[D](i)), nil
	})
}

// Lazy2 is Lazy for two-dependency constructors.
func Lazy2[T, D1, D2 any](ctor func(D1, D2) T) func(do.Injector) {
	return do.Lazy(func(i do.Injector) (T, error) {
		return ctor(do.MustInvoke[D1](i), do.MustInvoke[D2](i)), nil
	})
}

// Lazy3 is Lazy for three-dependency constructors.
func Lazy3[T, D1, D2, D3 any](ctor func(D1, D2, D3) T) func(do.Injector) {
	return do.Lazy(func(i do.Injector) (T, error) {
		return ctor(do.MustInvoke[D1](i), do.MustInvoke[D2](i), do.MustInvoke[D3](i)), nil
	})
}
