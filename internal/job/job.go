// Package job holds scheduled jobs, contributed under the "jobs:" naming
// convention and registered on the cron scheduler at startup.
package job

import (
	"context"
	"log/slog"

	"github.com/libtnb/cron"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/registry"
)

// Prefix marks job contributions; assemblers and tests collect by it.
const Prefix = "jobs:"

// Package lists every job contribution; add yours here.
var Package = do.Package(
	do.LazyNamed(Prefix+"heartbeat", Heartbeat),
)

// JobFn is a module's contribution to the scheduler. Specs accept an
// optional seconds field, @every descriptors and TZ= prefixes.
type JobFn func(c *cron.Cron) error

// Register applies every "jobs:*" contribution to c.
func Register(i do.Injector, c *cron.Cron) error {
	jobs, err := registry.Collect[JobFn](i, Prefix)
	if err != nil {
		return err
	}
	for _, register := range jobs {
		if err := register(c); err != nil {
			return err
		}
	}

	return nil
}

// Heartbeat is an example job; replace it with real ones.
func Heartbeat(i do.Injector) (JobFn, error) {
	log := do.MustInvoke[*slog.Logger](i)

	return func(c *cron.Cron) error {
		_, err := c.Add("@hourly", cron.JobFunc(func(ctx context.Context) error {
			log.InfoContext(ctx, "cron heartbeat")
			return nil
		}), cron.WithName("heartbeat"))
		return err
	}, nil
}
