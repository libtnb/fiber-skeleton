// Package registry collects same-kind contributions (routes, commands, jobs)
// registered under the naming conventions below — do's stand-in for value
// groups.
package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samber/do/v2"
)

// Contribution prefixes.
const (
	RoutePrefix      = "routes:"
	CommandPrefix    = "commands:"
	JobPrefix        = "jobs:"
	SubscriberPrefix = "subscribers:"
)

// Collect resolves every service under prefix, sorted for determinism.
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

// Verify fails on contributions whose prefix matches none of the known ones —
// a typo like "route:user" would otherwise be dropped silently.
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

// Lazy adapts a container-free one-dependency constructor into a provider.
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
