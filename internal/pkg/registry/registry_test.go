package registry_test

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

func TestCollectGathersByPrefixInNameOrder(t *testing.T) {
	i := do.New()
	do.ProvideNamedValue(i, registry.RoutePrefix+"b", "second")
	do.ProvideNamedValue(i, registry.RoutePrefix+"a", "first")
	do.ProvideNamedValue(i, registry.CommandPrefix+"x", "other kind")

	got, err := registry.Collect[string](i, registry.RoutePrefix)
	require.NoError(t, err)
	require.Equal(t, []string{"first", "second"}, got)
}

func TestCollectReportsWrongType(t *testing.T) {
	i := do.New()
	do.ProvideNamedValue(i, registry.RoutePrefix+"a", 42)

	_, err := registry.Collect[string](i, registry.RoutePrefix)
	require.Error(t, err)
}

func TestVerifyCatchesUnknownPrefix(t *testing.T) {
	i := do.New()
	do.ProvideNamedValue(i, registry.RoutePrefix+"ok", "fine")
	do.ProvideNamedValue(i, "route:user", "typo")

	require.NoError(t, registry.Verify(do.New(), registry.RoutePrefix))
	err := registry.Verify(i, registry.RoutePrefix, registry.CommandPrefix)
	require.ErrorContains(t, err, "route:user")
}

func TestVerifyIgnoresUnnamedServices(t *testing.T) {
	i := do.New()
	do.ProvideValue(i, 42) // type-derived name, no colon prefix

	require.NoError(t, registry.Verify(i, registry.RoutePrefix))
}

type widget struct{ deps int }

func TestLazyAdaptersInjectInOrder(t *testing.T) {
	i := do.New()
	do.ProvideValue(i, "dep")
	do.ProvideValue(i, 7)

	registry.Lazy(func(s string) *widget { return &widget{deps: 1} })(i)
	got, err := do.Invoke[*widget](i)
	require.NoError(t, err)
	require.Equal(t, 1, got.deps)

	j := do.New()
	do.ProvideValue(j, "dep")
	do.ProvideValue(j, 7)
	registry.Lazy2(func(s string, n int) widget { return widget{deps: 2} })(j)
	got2, err := do.Invoke[widget](j)
	require.NoError(t, err)
	require.Equal(t, 2, got2.deps)

	k := do.New()
	do.ProvideValue(k, "dep")
	do.ProvideValue(k, 7)
	do.ProvideValue(k, true)
	registry.Lazy3(func(s string, n int, b bool) uint { return 3 })(k)
	got3, err := do.Invoke[uint](k)
	require.NoError(t, err)
	require.EqualValues(t, 3, got3)
}
