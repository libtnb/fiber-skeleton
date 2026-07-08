package bootstrap

import (
	"log/slog"

	"github.com/libtnb/cron"
	"github.com/libtnb/cron/wrap"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/job"
)

// NewCron builds the scheduler and registers every job contribution on it.
func NewCron(i do.Injector) (*cron.Cron, error) {
	c := cron.New(
		cron.WithLogger(do.MustInvoke[*slog.Logger](i)),
		cron.WithSecondsField(),
		cron.WithChain(wrap.SkipIfRunning()),
	)

	if err := job.Register(i, c); err != nil {
		return nil, err
	}

	return c, nil
}
