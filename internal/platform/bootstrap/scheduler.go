package bootstrap

import (
	"context"
	"log/slog"

	"github.com/libtnb/cron"
	"github.com/libtnb/cron/wrap"

	"github.com/libtnb/fiber-skeleton/internal/shared/job"
	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
)

func NewCron(log *slog.Logger, jobs registry.Jobs) (*cron.Cron, error) {
	c, err := cron.New(
		cron.WithLogger(log),
		cron.WithSecondsField(),
		cron.WithChain(wrap.SkipIfRunning()),
	)
	if err != nil {
		return nil, err
	}
	if err := registerJobs(jobs, c); err != nil {
		return nil, err
	}

	return c, nil
}

func registerJobs(jobs registry.Jobs, c *cron.Cron) error {
	for _, apply := range jobs {
		if err := apply(c); err != nil {
			return err
		}
	}

	return nil
}

// Heartbeat is an example job; replace it with real ones.
func Heartbeat(log *slog.Logger) job.Fn {
	return func(c *cron.Cron) error {
		_, err := c.Add("@hourly", cron.JobFunc(func(ctx context.Context) error {
			log.InfoContext(ctx, "cron heartbeat")
			return nil
		}), cron.WithName("heartbeat"))
		return err
	}
}
