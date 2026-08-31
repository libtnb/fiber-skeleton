// Package job is the scheduler contribution contract shared by every module.
package job

import "github.com/libtnb/cron"

// Fn is a module's scheduler contribution assembled by wire.
type Fn func(c *cron.Cron) error
