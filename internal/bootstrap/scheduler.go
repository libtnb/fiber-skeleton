package bootstrap

import (
	"context"
	"log/slog"

	"github.com/libtnb/cron"
	"github.com/libtnb/cron/wrap"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/job"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

func NewCron(i do.Injector) (*cron.Cron, error) {
	c := cron.New(
		cron.WithLogger(do.MustInvoke[*slog.Logger](i)),
		cron.WithSecondsField(),
		cron.WithChain(wrap.SkipIfRunning()),
	)

	if err := registerJobs(i, c); err != nil {
		return nil, err
	}

	return c, nil
}

func registerJobs(i do.Injector, c *cron.Cron) error {
	jobs, err := registry.Collect[job.Fn](i, registry.JobPrefix)
	if err != nil {
		return err
	}
	for _, apply := range jobs {
		if err := apply(c); err != nil {
			return err
		}
	}

	return nil
}

// Heartbeat is an example job; replace it with real ones.
func Heartbeat(i do.Injector) (job.Fn, error) {
	log := do.MustInvoke[*slog.Logger](i)

	return func(c *cron.Cron) error {
		_, err := c.Add("@hourly", cron.JobFunc(func(ctx context.Context) error {
			log.InfoContext(ctx, "cron heartbeat")
			return nil
		}), cron.WithName("heartbeat"))
		return err
	}, nil
}
